package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/daf/internal/viewmodel"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/fpawel/serial"
	"github.com/hako/durafmt"
	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
	"time"
)

func isFailWork(err error) bool {
	return err != nil && !isDeviceError(err)
}

func isDeviceError(err error) bool {
	return merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded)
}

type Product = viewmodel.DafProductViewModel

func doForEachOkProduct(work func(p *Product) error) error {
	if len(prodsMdl.OkProducts()) == 0 {
		return ErrNoOkProducts
	}
	for _, p := range prodsMdl.OkProducts() {
		if err := work(p); isFailWork(err) {
			return err
		}
	}
	return nil
}

func dafSendCmdToEachOkProduct(cmd modbus.DevCmd, arg float64) error {

	dafMainWindow.SetWorkStatus(walk.RGB(0, 0, 128),
		fmt.Sprintf("отправка команды %X, %v", cmd, arg))

	if cmd == 5 {
		if err := portDaf.open(); err != nil {
			return err
		}
		_, err := portDaf.Port.Write(modbus.Write32BCDRequest(0, 0x10, cmd, arg).Bytes())
		return err
	}

	return doForEachOkProduct(func(p *Product) error {
		return dafSendCmdToPlace(p.Place, cmd, arg)
	})
}

func blowGas(gas data.Gas) error {
	if err := switchGas(gas); err != nil {
		if !dafMainWindow.IgnoreErrorPrompt(fmt.Sprintf("газовый блок: %d", gas), err) {
			return err
		}
	}
	t := 5 * time.Minute
	if gas == 1 {
		t = 10 * time.Minute
	}
	return delay(fmt.Sprintf("продувка ПГС%d", gas), t)
}

func dafSetupCurrent() error {
	setCurrentWorkName("настройка токового выхода")

	if err := dafSendCmdToEachOkProduct(0xB, 1); err != nil {
		return err
	}

	sleep(5 * time.Second)

	dafMainWindow.SetWorkStatus(ColorNavy, "корректировка тока 4 мА")

	if err := doForEachOkProduct(func(p *Product) error {
		v, err := EN6408Read(p.Place) // считать тока со стенда
		if err != nil {
			return err
		}
		return dafSendCmdToPlace(p.Place, 9, v.Current)
	}); err != nil {
		return err
	}

	if err := dafSendCmdToEachOkProduct(0xC, 1); err != nil {
		return err
	}

	sleep(5 * time.Second)

	dafMainWindow.SetWorkStatus(ColorNavy, "корректировка тока 20 мА")
	if err := doForEachOkProduct(func(p *Product) error {
		v, err := EN6408Read(p.Place)
		if err != nil {
			return err
		}
		return dafSendCmdToPlace(p.Place, 0xA, v.Current)
	}); err != nil {
		return err
	}
	return nil
}

func dafSetupThresholdTest() error {

	setCurrentWorkName("установка порогов для настройки")

	party := data.GetLastParty()

	if err := dafSendCmdToEachOkProduct(0x30, party.Threshold1Test); err != nil {
		return err
	}
	if err := dafSendCmdToEachOkProduct(0x31, party.Threshold2Test); err != nil {
		return err
	}

	return nil
}

func dafAdjust() error {

	party := data.GetLastParty()

	setCurrentWorkName("корректировка нулевых показаний")

	defer func() {
		_ = switchGas(0)
	}()

	if err := blowGas(1); err != nil {
		return err
	}
	if err := dafSendCmdToEachOkProduct(0x32, party.Pgs1); err != nil {
		return err
	}

	setCurrentWorkName("корректировка чувствительности")

	if err := blowGas(4); err != nil {
		return err
	}
	if err := dafSendCmdToEachOkProduct(0x33, party.Pgs4); err != nil {
		return err
	}
	return nil
}

func dafTestMeasureRange() error {

	defer func() {
		_ = switchGas(0)
	}()

	setCurrentWorkName("проверка диапазона измерений")

	for n, gas := range []data.Gas{1, 2, 3, 4, 3, 1} {

		what := fmt.Sprintf("проверка диапазона измерений: ПГС%d, точка %d", gas, n+1)

		dafMainWindow.SetWorkStatus(ColorNavy, what+": продувка газа")

		if err := blowGas(gas); err != nil {
			return err
		}

		dafMainWindow.SetWorkStatus(ColorNavy, what+": опрос и сохранение данных")

		if err := doForEachOkProduct(func(p *Product) error {

			dv, err := dafReadAtPlace(p.Place)
			if isFailWork(err) {
				return nil
			}
			v, err := EN6408Read(p.Place)
			if err != nil {
				return nil
			}
			value := data.ProductValue{
				ProductID:     p.ProductID,
				Gas:           gas,
				CreatedAt:     time.Now(),
				WorkIndex:     n,
				Concentration: dv.Concentration,
				Current:       v.Current,
				Threshold1:    v.Threshold1,
				Threshold2:    v.Threshold2,
				Mode:          dv.Mode,
				FailureCode:   dv.Failure,
			}

			data.DBxProducts.MustExec(
				`DELETE FROM product_value WHERE product_id = ? AND name = ? AND work_index = ?`,
				p.ProductID, "диапазон измерений", n)

			if err := data.DBProducts.Save(&value); err != nil {
				panic(err)
			}
			logrus.Infof("сохранено значение: место %d, адрес %d: %v", p.Place, p.Addr, value)

			return nil

		}); err != nil {
			return err
		}
	}

	if err := blowGas(1); err != nil {
		return err
	}

	return nil
}

func dafTestStability() error {
	if err := blowGas(3); err != nil {
		return err
	}
	return nil
}

func dafSetupMain() error {
	for _, f := range []func() error{dafSetupCurrent, dafSetupThresholdTest, dafAdjust, dafTestMeasureRange} {
		if err := f(); err != nil {
			return err
		}
	}
	if data.GetLastParty().Type == 0 {
		return nil
	}

	setCurrentWorkName("проверка HART протокола")
	return doForEachOkProduct(testHart)
}

func interrogateProducts() error {
	currentWorkName = ""
	for {
		if len(prodsMdl.OkProducts()) == 0 {
			return errors.New("не выбрано ни одной строки в таблице приборов текущей партии")
		}
		for _, p := range prodsMdl.OkProducts() {
			if _, err := EN6408Read(p.Place); err != nil {
				return err
			}
			if _, err := dafReadAtPlace(p.Place); isFailWork(err) {
				return err
			}
		}
	}
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
	logrus.Infoln("пауза", durafmt.Parse(t))
	defer func() {
		logrus.Infoln("окончание паузы", durafmt.Parse(t))
	}()
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
		dafMainWindow.DelayHelp.Hide()
	}()

	startMoment := time.Now()

	skipDelay = func() {
		doSkipDelay()
		logrus.Warnf("%s %s: задержка прервана: %s", what, durafmt.Parse(total),
			durafmt.Parse(time.Since(startMoment)))
	}
	dafMainWindow.DelayHelp.Show(what, total)

	for {
		for _, p := range prodsMdl.OkProducts() {
			_, err := dafReadAtPlace(p.Place)
			if ctxDelay.Err() != nil {
				return nil
			}
			if isFailWork(err) {
				return err
			}
		}

		func() {
			timer := time.NewTimer(5 * time.Second)
			for {
				select {
				case <-timer.C:
					return
				case <-ctxDelay.Done():
					timer.Stop()
					return
				}
			}
		}()

	}
}

func onPlaceConnectionError(place int, err error) {
	p := prodsMdl.ProductAt(place)
	logrus.Errorf("место %d, адрес %d, серийный номер %d, ID %d: %v",
		place+1, p.Addr, p.Serial, p.ProductID, err)
	prodsMdl.SetConnectionErrorAt(place, err)
	if currentWorkName != "" {
		data.WriteProductError(p.ProductID, currentWorkName, err)
	}
}

func setCurrentWorkName(workName string) {
	data.DBxProducts.MustExec(
		`DELETE FROM product_entry 
WHERE work_name = ? 
  AND product_id IN (SELECT product_id FROM last_party_products)`, workName)
	currentWorkName = workName
}

var (
	dafMainWindow   DafMainWindow
	prodsMdl        *viewmodel.DafProductsTable
	currentWorkName string

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

	ErrNoOkProducts = merry.New("отстутсвуют приборы, которые отмеченны галочками и не имеют ошибок связи")
)
