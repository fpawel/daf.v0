package main

import (
	"context"
	"errors"
	"github.com/ansel1/merry"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/daf/internal/viewmodel"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/fpawel/serial"
	"github.com/hako/durafmt"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	WorkSetupCurrent = "настройка токового выхода"
)

func setupCurrents() error {
	data.ClearCurrentProductsWorkResult(WorkSetupCurrent)
	defer func() {

	}()

	for place, p := range prodsMdl.OkProducts() {
		if err := sendCmdPlace(place, p.Addr, 0xB, 1); err != nil && !IsDeviceError(err) {
			return err
		}
	}

	sleep(5 * time.Second)

	for place, p := range data.GetProductsOfLastParty() {
		v, err := read6408(place, p.Addr)
		if err != nil {
			return err
		}
		if err = sendCmdPlace(place, p.Addr, 9, v.Current); err != nil && !IsDeviceError(err) {
			return err
		}
		if err = sendCmdPlace(place, p.Addr, 0xC, 1); err != nil && !IsDeviceError(err) {
			return err
		}
	}

	sleep(5 * time.Second)

	for place, p := range prodsMdl.OkProducts() {
		v, err := read6408(place, p.Addr)
		if err != nil {
			return err
		}
		if err := sendCmdPlace(place, p.Addr, 0xA, v.Current); err != nil && !IsDeviceError(err) {
			return err
		}
		data.SetProductWorkInfo(p.ProductID, WorkSetupCurrent, "выполнено")
	}

	return nil
}

func interrogateProducts() error {
	for {
		if len(prodsMdl.OkProducts()) == 0 {
			return errors.New("не выбрано ни одной строки в таблице приборов текущей партии")
		}
		for place, p := range prodsMdl.OkProducts() {
			if _, err := read6408(place, p.Addr); err != nil {
				return err
			}
			if _, err := readDaf(place, p.Addr); err != nil && !IsDeviceError(err) {
				return err
			}
		}
	}
}

func read6408(place int, addr modbus.Addr) (*viewmodel.DafValue6408, error) {

	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()
	b, err := modbus.Read3(portDaf, 32, modbus.Var(addr-1)*2, 2, func(_, _ []byte) error {
		return nil
	})
	if err != nil {
		return nil, merry.Append(err, "запрос тока и состояния контактов реле")
	}
	b = b[3:]
	v := new(viewmodel.DafValue6408)
	v.Current = (float64(b[0])*256 + float64(b[1])) / 100
	v.Threshold1 = b[3]&1 == 0
	v.Threshold2 = b[3]&2 == 0
	logrus.Debugf("адрес %d: %v", addr, *v)
	prodsMdl.Set6408Value(place, *v)
	return v, nil
}

func readDaf(place int, addr modbus.Addr) (v viewmodel.DafValue, err error) {
	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()
	v, err = doReadDaf(addr)
	if err != nil {
		if merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded) {
			prodsMdl.SetConnectionErrorAt(place, err)
		}
		logrus.Errorf("место %d, адрес %d: %v", place+1, addr, err)
		return
	}
	logrus.Debugf("место %d, адрес %d: %v", place+1, addr, v)
	prodsMdl.SetDafValue(place, v)
	return
}

func doReadDaf(addr modbus.Addr) (v viewmodel.DafValue, err error) {

	for _, x := range []struct {
		var3 modbus.Var
		p    *float64
	}{
		{0x00, &v.Concentration},
		{0x1C, &v.Threshold1},
		{0x1E, &v.Threshold2},
		{0x20, &v.Failure},
		{0x36, &v.Version},
		{0x3A, &v.VersionID},
		{0x32, &v.Gas},
	} {
		if *x.p, err = modbus.Read3BCD(portDaf, addr, x.var3); err != nil {
			return
		}
	}
	v.Mode, err = modbus.ReadUInt16(portDaf, addr, 0x23)
	return
}

func sendCmd(cmd modbus.DevCmd, arg float64) error {

	if cmd == 5 {
		if err := portDaf.open(); err != nil {
			return err
		}
		_, err := portDaf.Port.Write(modbus.Write32BCDRequest(0, 0x10, cmd, arg).Bytes())
		return err
	}

	for place, p := range prodsMdl.OkProducts() {
		if err := sendCmdPlace(place, p.Addr, cmd, arg); err != nil && !IsDeviceError(err) {
			return err
		}
	}
	return nil
}

func sendCmdPlace(place int, addr modbus.Addr, cmd modbus.DevCmd, arg float64) error {
	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()

	err := modbus.Write32FloatProto(portDaf, addr, 0x10, cmd, arg)
	if err == nil {
		logrus.Infof("ДАФ №%d, адрес %d: запись в 32-ой регистр %X, %v", place+1, addr, cmd, arg)
		prodsMdl.SetConnectionOkAt(place)
		return nil
	}
	logrus.Errorf("ДАФ №%d, адрес %d: %v", place+1, addr, err)
	prodsMdl.SetConnectionErrorAt(place, err)
	return err
}

func switchGas(gas data.Gas) error {

	req := modbus.Req{
		Addr:     33,
		ProtoCmd: 0x10,
		Data:     []byte{0, 32, 0, 1, 2, 0, byte(gas)},
	}
	_, err := portDaf.GetResponse(req.Bytes(), func(_, response []byte) error {
		return req.CheckResponse(response)
	})
	if err != nil {
		err = merry.Appendf(err, "газовый блок: %d", gas)
	}
	return err
}

func IsDeviceError(err error) bool {
	return merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded)
}

type port struct {
	Port     *comport.Port
	PortName string
	Config   comm.Config
}

func (x port) open() error {
	if x.Port.Opened() {
		return nil
	}
	return x.Port.Open(x.PortName)
}

func (x port) GetResponse(request []byte, prs comm.ResponseParser) ([]byte, error) {
	if err := x.open(); err != nil {
		return nil, err
	}
	return x.Port.GetResponse(request, x.Config, comportContext, prs)
}

func sleep(t time.Duration) {
	timer := time.NewTimer(t)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return
		case <-comportContext.Done():
			return
		}
	}
}

func delay(what string, total time.Duration) error {

	originalComportContext := comportContext
	ctxDelay, doSkipDelay := context.WithTimeout(comportContext, total)
	comportContext = ctxDelay
	defer func() {
		comportContext = originalComportContext
		guiHideDelay()
	}()
	startMoment := time.Now()

	skipDelay = func() {
		doSkipDelay()
		logrus.Warnf("%s %s: задержка прервана: %s", what, durafmt.Parse(total),
			durafmt.Parse(time.Since(startMoment)))
	}
	guiShowDelay(what, total)

	for {
		if len(prodsMdl.OkProducts()) == 0 {
			return errors.New("не выбрано ни одной строки в таблице приборов текущей партии")
		}
		for place, p := range prodsMdl.OkProducts() {
			if ctxDelay.Err() != nil {
				return nil
			}
			if !p.Checked {
				continue
			}

			_, err := readDaf(place, p.Addr)
			if err == nil ||
				merry.Is(err, comm.ErrProtocol) ||
				merry.Is(err, context.DeadlineExceeded) {
				continue
			}
			if merry.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}

var (
	prodsMdl *viewmodel.DafProductsTable

	portDaf = port{
		Port: comport.NewPort("стенд", serial.Config{Baud: 9600}, func(entry comport.Entry) {
			logrus.Debugln(entry)
		}),
		Config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       2,
		},
	}

	portHart = port{
		Port: comport.NewPort("hart", serial.Config{
			Baud:        1200,
			ReadTimeout: time.Millisecond,
			Parity:      serial.ParityOdd,
			StopBits:    serial.Stop1,
		}, func(entry comport.Entry) {
			logrus.Debugln(entry)
		}),
		Config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     2000,
			MaxAttemptsRead:       5,
		},
	}

	cancelComport  = func() {}
	skipDelay      = func() {}
	comportContext context.Context

	guiShowDelay func(what string, total time.Duration)
	guiHideDelay func()
)
