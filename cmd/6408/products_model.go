package main

import (
	"fmt"
	"github.com/fpawel/daf/internal/assets"
	"github.com/fpawel/daf/internal/data"
	"github.com/lxn/walk"
)

type productsModel struct {
	walk.TableModelBase
	items       []productModel
	interrogate int
}

type productModel struct {
	*data.Product
	Concentration,
	Current *float64
	Threshold1,
	Threshold2 *bool
}

func (m *productsModel) addNewProduct() {
	serial := int64(1)
	addr := 1
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

func (m *productsModel) validate() {
	m.interrogate = -1
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
		if x.Concentration != nil {
			return *x.Concentration
		}
	case 3:
		if x.Current != nil {
			return *x.Current
		}
	}
	return ""
}

func (m *productsModel) StyleCell(style *walk.CellStyle) {
	p := m.items[style.Row()]
	switch style.Col() {
	case 0:
		if style.Row() == m.interrogate {
			style.Image = assets.ImgForward
		}
	case 4:
		if p.Threshold1 != nil {
			if *p.Threshold1 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}

	case 5:
		if p.Threshold2 != nil {
			if *p.Threshold2 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
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
