package main

import (
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/daf/internal/viewmodel"
	"github.com/powerman/structlog"
	"github.com/sirupsen/logrus"
)

type EN6408ConnectionLine int

const (
	EN6408Disconnect EN6408ConnectionLine = iota
	EN6408ConnectRS485
	EN6408ConnectHart
)

var (
	ErrEN6408   = merry.New("стенд 6408")
	ErrGasBlock = merry.New("газовый блок")
)

func EN6408SetConnectionLine(place int, connLine EN6408ConnectionLine) error {
	req := modbus.Req{
		Addr:     0x20,
		ProtoCmd: 0x10,
		Data:     []byte{0, byte(place), 0, 1, 2, 0, byte(connLine)},
	}
	_, err := req.GetResponse(portDaf, nil)
	if err != nil {
		return ErrEN6408.Appendf("место %d: выбор линии связи %d: %v", place+1, connLine, err)
	}
	return err
}

func EN6408Read(place int) (*viewmodel.DafValue6408, error) {

	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()

	addr := prodsMdl.ProductAt(place).Addr
	b, err := modbus.Read3(portDaf, 32, modbus.Var(addr-1)*2, 2, func(_, _ []byte) error {
		return nil
	})
	if err != nil {
		return nil, ErrEN6408.Appendf("опрос места %d: %+v", place+1, err)
	}
	b = b[3:]
	v := new(viewmodel.DafValue6408)
	v.Current = (float64(b[0])*256 + float64(b[1])) / 100
	v.Threshold1 = b[3]&1 == 0
	v.Threshold2 = b[3]&2 == 0
	prodsMdl.Set6408Value(place, *v)

	structlog.DefaultLogger.Info("ЭН6408: опрос места", KeyPlace, place, KeyEN6408, *v)

	return v, nil
}

func switchGas(gas data.Gas) error {

	req := modbus.Req{
		Addr:     33,
		ProtoCmd: 0x10,
		Data:     []byte{0, 32, 0, 1, 2, 0, byte(gas)},
	}
	if _, err := req.GetResponse(portDaf, nil); err != nil {
		return ErrGasBlock.Appendf("клапан %d: %v", gas, err)
	}
	return nil
}

func dafReadAtPlace(place int) (v viewmodel.DafValue, err error) {

	addr := prodsMdl.ProductAt(place).Addr

	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()

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
			break
		}
	}
	if err == nil {
		v.Mode, err = modbus.ReadUInt16(portDaf, addr, 0x23)
	}

	if err == nil {
		logrus.Debugf("место %d, адрес %d: %v", place+1, addr, v)

		log.Debug(fmt.Sprintf("%+v", v),
			structlog.KeyTime, now(),
			"место", place+1,
			"адрес", addr,
		)

		prodsMdl.SetDafValue(place, v)
	}
	if isDeviceError(err) {
		onPlaceConnectionError(place, err)
	}
	return
}

func dafSendCmdToPlace(place int, cmd modbus.DevCmd, arg float64) error {
	prodsMdl.SetInterrogatePlace(place)
	defer func() {
		prodsMdl.SetInterrogatePlace(-1)
	}()

	addr := prodsMdl.ProductAt(place).Addr

	err := modbus.Write32FloatProto(portDaf, addr, 0x10, cmd, arg)
	if err == nil {
		logrus.Infof("ДАФ №%d, адрес %d: запись в 32-ой регистр %X, %v", place+1, addr, cmd, arg)
		prodsMdl.SetConnectionOkAt(place)
		return nil
	}
	if isDeviceError(err) {
		onPlaceConnectionError(place, err)
	}
	return err
}
