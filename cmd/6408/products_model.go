package main

import (
	"fmt"
	"github.com/fpawel/daf/internal/assets"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/lxn/walk"
)

type productsModel struct {
	walk.TableModelBase
	items            []productModel
	interrogatePlace int
}

type productModel struct {
	*data.Product
	Value6408  *Value6408
	Daf        *DafValue
	connection *connectionInfo
}

type Value6408 struct {
	Current                float64
	Threshold1, Threshold2 bool
}

type DafValue struct {
	Concentration,
	Threshold1, Threshold2,
	Failure, Version, VersionID, Gas float64

	Mode uint16
}

type connectionInfo struct {
	ok   bool
	text string
}

type ProductColumn int

const (
	pcAddr ProductColumn = iota
	pcSerialNumber
	pcProductID
	pcConcentration
	pcCurrent
	pcMode
	pcFailure
	pcThreshold1
	pcThreshold2
	pcVersion
	pcGas
	pcConnection
)

func (x Value6408) String() string {
	f := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	return fmt.Sprintf("ток=%v П1=%d П2=%d", x.Current, f(x.Threshold1), f(x.Threshold2))
}

func (x DafValue) String() string {
	return fmt.Sprintf("концентрация=%v режим=%v отказ=%v версия=%v порог1=%v порог2=%v",
		x.Concentration, x.Mode, x.Failure, x.Version, x.Threshold1, x.Threshold2)
}

func (m *productsModel) addNewProduct() {
	serial := int64(1)
	addr := modbus.Addr(1)
l1:
	for _, p := range m.items {
		if p.Addr == addr {
			addr++
			goto l1
		}
		if p.Serial == serial {
			serial++
			goto l1
		}
	}

	data.DBProducts.MustExec(
		`INSERT INTO product (party_id, serial, addr, checked) VALUES ((SELECT party_id FROM last_party), ?, ?, 1)`,
		serial, addr)
	m.validate()
}

func (m *productsModel) setDafValue(place int, v DafValue) {
	m.items[place].Daf = &v
	m.items[place].connection = &connectionInfo{true, v.String()}
	mainWindow.Synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *productsModel) set6408Value(place int, v Value6408) {
	m.items[place].Value6408 = &v
	mainWindow.Synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *productsModel) setProductConnection(place int, ok bool, text string) {
	m.items[place].connection = &connectionInfo{ok, text}
	mainWindow.Synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *productsModel) setInterrogatePlace(place int) {
	if m.interrogatePlace == place {
		return
	}
	m.interrogatePlace = place
	mainWindow.Synchronize(func() {
		m.PublishRowsReset()
	})
}

func (m *productsModel) validate() {
	m.interrogatePlace = -1
	m.items = nil
	for _, p := range data.GetProductsOfLastParty() {
		m.items = append(m.items, productModel{Product: p})
	}
	m.PublishRowsReset()
}

func (m *productsModel) RowCount() int {
	return len(m.items)
}

func (m *productsModel) Value(row, col int) interface{} {
	x := m.items[row]

	switch ProductColumn(col) {
	case pcAddr:
		return x.Addr
	case pcSerialNumber:
		return x.Serial
	case pcProductID:
		return x.ProductID
	case pcConnection:
		if x.connection != nil {
			return x.connection.text
		}
	case pcCurrent:
		if x.Value6408 != nil {
			return x.Value6408.Current
		}
	default:
		if x.Daf != nil {
			switch ProductColumn(col) {
			case pcConcentration:
				return x.Daf.Concentration
			case pcThreshold1:
				return x.Daf.Threshold1
			case pcThreshold2:
				return x.Daf.Threshold2
			case pcMode:
				return x.Daf.Mode
			case pcFailure:
				return int(x.Daf.Failure)
			case pcVersion:
				return fmt.Sprintf("%v.%X", x.Daf.Version, int(x.Daf.VersionID))
			case pcGas:
				return int(x.Daf.Gas)
			}
		}
	}
	return ""
}

func (m *productsModel) StyleCell(style *walk.CellStyle) {

	if style.Row() == m.interrogatePlace {
		style.BackgroundColor = walk.RGB(166, 202, 240)
	}

	p := m.items[style.Row()]
	switch ProductColumn(style.Col()) {
	case pcAddr:
		if style.Row() == m.interrogatePlace {
			style.Image = assets.ImgForward
		}
	case pcThreshold1:
		if p.Value6408 != nil {
			if p.Value6408.Threshold1 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}

	case pcThreshold2:
		if p.Value6408 != nil {
			if p.Value6408.Threshold2 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}
	case pcConnection:

		if p.connection != nil {
			if p.connection.ok {
				style.TextColor = walk.RGB(0, 0, 128)
				style.Image = assets.ImgCheckMark
			} else {
				style.TextColor = walk.RGB(255, 0, 0)
				style.Image = assets.ImgError
			}
		}

	}
}

func (m *productsModel) Checked(index int) bool {
	return m.items[index].Checked
}

func (m *productsModel) SetChecked(index int, checked bool) error {
	m.items[index].Checked = checked
	_, err := data.DBProducts.Exec(`UPDATE product SET checked = ? WHERE product_id = ?`, checked, m.items[index].ProductID)
	return err
}

var (
	lastPartyProductsModel = &productsModel{}
)

func init() {
	lastPartyProductsModel.validate()
}
