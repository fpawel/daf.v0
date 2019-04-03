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
	Value      *ProductValue
	connection *connectionInfo
}

type ProductValue struct {
	Concentration, Current float64
	Threshold1, Threshold2 bool
}

type connectionInfo struct {
	ok   bool
	text string
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

func (m *productsModel) setProductValue(place int, v ProductValue) {
	m.items[place].Value = &v
	m.items[place].connection = &connectionInfo{true, "установлена"}
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
	switch col {
	case 0:
		return x.Addr
	case 1:
		return fmt.Sprintf("%d-%d", x.Serial, x.ProductID)
	case 2:
		if x.Value != nil {
			return x.Value.Concentration
		}
	case 3:
		if x.Value != nil {
			return x.Value.Current
		}
	case 6:
		if x.connection != nil {
			return x.connection.text
		}
	}
	return ""
}

func (m *productsModel) StyleCell(style *walk.CellStyle) {

	if style.Row() == m.interrogatePlace {
		style.BackgroundColor = walk.RGB(166, 202, 240)
	}

	p := m.items[style.Row()]
	switch style.Col() {
	case 0:
		if style.Row() == m.interrogatePlace {
			style.Image = assets.ImgForward
		}
	case 1:
		if p.connection != nil {
			if p.connection.ok {
				style.Image = assets.ImgCheckMark
			} else {
				style.Image = assets.ImgError
			}
		}
	case 4:
		if p.Value != nil {
			if p.Value.Threshold1 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}

	case 5:
		if p.Value != nil {
			if p.Value.Threshold2 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}
	case 6:
		if p.connection != nil {
			if p.connection.ok {
				style.TextColor = walk.RGB(0, 0, 128)
			} else {
				style.TextColor = walk.RGB(255, 0, 0)
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

var lastPartyProductsModel = &productsModel{}

func init() {
	lastPartyProductsModel.validate()
}
