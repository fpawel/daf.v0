package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/fpawel/serial"
	"github.com/sirupsen/logrus"
	"time"
)

func interrogateProducts() error {
	for {
		if !data.LastPartyHasCheckedProduct() {
			return errors.New("не выбрано ни одной строки в таблице приборов текущей партии")
		}
		for place, p := range data.GetProductsOfLastParty() {
			if !p.Checked {
				continue
			}
			_, err := readProduct(place, p.Addr)
			if err == nil || merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded) {
				continue
			}
			if merry.Is(err, context.Canceled) {
				return nil
			}
			return err
		}
	}
}

func readProduct(place int, addr modbus.Addr) (ProductValue, error) {
	lastPartyProductsModel.setInterrogatePlace(place)
	v, err := doReadProduct(addr)
	if err == nil {
		lastPartyProductsModel.setProductValue(place, v)
	} else if merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded) {
		lastPartyProductsModel.setProductConnection(place, false, err.Error())
	}

	if err == nil {
		logrus.Infof("№%d-ADDR%d: конц=%v ток=%v П1=%v П2=%v", place+1, addr,
			v.Concentration, v.Current, v.Threshold1, v.Threshold2)
	} else {
		logrus.Errorf("№%d-ADDR%d: %v", place+1, addr, err)
	}

	return v, err
}

func doReadProduct(addr modbus.Addr) (v ProductValue, err error) {

	if v.Concentration, err = modbus.Read3BCD(portDaf, addr, 0); err != nil {
		return
	}

	var b []byte
	b, err = modbus.Read3(portDaf, 32, modbus.Var(addr-1)*2, 2, func(_, _ []byte) error {
		return nil
	})
	if err != nil {
		return
	}
	b = b[3:]

	v.Current = (float64(b[0])*256 + float64(b[1])) / 100
	v.Threshold1 = b[3]&1 == 0
	v.Threshold2 = b[3]&2 == 0

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

	for place, p := range data.GetProductsOfLastParty() {
		if !p.Checked {
			continue
		}
		lastPartyProductsModel.setInterrogatePlace(place)

		err := modbus.Write32FloatProto(portDaf, p.Addr, 0x10, cmd, arg)

		what := fmt.Sprintf("№%d-ADDR%d-SER%d-ID%d: команда %d, %v", place+1, p.Addr, p.Serial, p.ProductID, cmd, arg)
		if err == nil {
			logrus.Info(what)
		} else {
			logrus.Errorf("%s: %v", what, err)
		}

		if err == nil {
			lastPartyProductsModel.setProductConnection(place, true, fmt.Sprintf("команда %d, %v", cmd, arg))
		} else {
			if merry.Is(err, context.Canceled) {
				return nil
			}
			if merry.Is(err, comm.ErrProtocol) || merry.Is(err, context.DeadlineExceeded) {
				lastPartyProductsModel.setProductConnection(place, false, fmt.Sprintf("%v: команда %d, %v", err, cmd, arg))
			} else {
				return err
			}
		}
	}
	return nil
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

func onComport(entry comport.Entry) {
	if entry.Error == nil {
		logrus.Debugln(entry)
	} else {
		logrus.Errorln(entry)
	}

}

var (
	portDaf = port{
		Port: comport.NewPort("стенд", serial.Config{Baud: 9600}, onComport),
		Config: comm.Config{
			ReadByteTimeoutMillis: 30,
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
		}, onComport),
		Config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     2000,
			MaxAttemptsRead:       5,
		},
	}

	cancelComport  = func() {}
	comportContext context.Context
)
