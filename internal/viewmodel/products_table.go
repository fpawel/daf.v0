package viewmodel

import (
	"fmt"
	"github.com/fpawel/daf/internal/assets"
	"github.com/fpawel/daf/internal/data"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

type synchronizeFunc = func(func())

type DafProductsTable struct {
	walk.TableModelBase
	synchronize      synchronizeFunc
	items            []DafProductViewModel
	interrogatePlace int
}

type DafProductViewModel struct {
	*data.Product
	Value6408  *DafValue6408
	Daf        *DafValue
	connection *connectionInfo
}

type DafValue6408 struct {
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

func (x DafValue6408) String() string {
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

func NewDafProductsTable(synchronize synchronizeFunc) *DafProductsTable {
	return &DafProductsTable{
		synchronize: synchronize,
	}
}

func (m *DafProductsTable) AddNewProduct() {
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

	if err := data.DBProducts.Save(&data.Product{
		PartyID:   data.GetLastPartyID(),
		Addr:      addr,
		Serial:    serial,
		Checked:   true,
		CreatedAt: time.Now(),
	}); err != nil {
		panic(err)
	}
	m.Validate()
}

func (m *DafProductsTable) SetDafValue(place int, v DafValue) {
	m.items[place].Daf = &v
	m.items[place].connection = &connectionInfo{true, v.String()}
	m.synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *DafProductsTable) Set6408Value(place int, v DafValue6408) {
	m.items[place].Value6408 = &v
	m.synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *DafProductsTable) SetProductConnection(place int, ok bool, text string) {
	m.items[place].connection = &connectionInfo{ok, text}
	m.synchronize(func() {
		m.PublishRowChanged(place)
	})
}

func (m *DafProductsTable) SetInterrogatePlace(place int) {
	if m.interrogatePlace == place {
		return
	}
	n := m.interrogatePlace
	m.interrogatePlace = place
	m.synchronize(func() {
		m.PublishRowChanged(n)
		m.PublishRowChanged(place)
	})
}

func (m *DafProductsTable) Validate() {
	m.interrogatePlace = -1
	m.items = nil
	for _, p := range data.GetProductsOfLastParty() {
		m.items = append(m.items, DafProductViewModel{Product: p})
	}
	m.synchronize(func() {
		m.PublishRowsReset()
	})
}

func (m *DafProductsTable) RowCount() int {
	return len(m.items)
}

func (m *DafProductsTable) ProductAt(n int) *data.Product {
	return m.items[n].Product
}

func (m *DafProductsTable) Value(row, col int) interface{} {
	x := m.items[row]

	switch ProductColumn(col) {
	case ProdColAddr:
		return x.Addr
	case ProdColSerialNumber:
		return x.Serial
	case ProdColProductID:
		return x.ProductID
	case ProdColConnection:
		if x.connection != nil {
			return x.connection.text
		}
	case ProdColCurrent:
		if x.Value6408 != nil {
			return x.Value6408.Current
		}
	default:
		if x.Daf != nil {
			switch ProductColumn(col) {
			case ProdColConcentration:
				return x.Daf.Concentration
			case ProdColThreshold1:
				return x.Daf.Threshold1
			case ProdColThreshold2:
				return x.Daf.Threshold2
			case ProdColMode:
				return x.Daf.Mode
			case ProdColFailure:
				return int(x.Daf.Failure)
			case ProdColVersion:
				return fmt.Sprintf("%v.%X", x.Daf.Version, int(x.Daf.VersionID))
			case ProdColGas:
				return int(x.Daf.Gas)
			}
		}
	}
	return ""
}

func (m *DafProductsTable) StyleCell(style *walk.CellStyle) {

	if style.Row() == m.interrogatePlace {
		style.BackgroundColor = walk.RGB(166, 202, 240)
	}

	p := m.items[style.Row()]
	switch ProductColumn(style.Col()) {
	case ProdColAddr:
		if style.Row() == m.interrogatePlace {
			style.Image = assets.ImgForward
		}
	case ProdColThreshold1:
		if p.Value6408 != nil {
			if p.Value6408.Threshold1 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}

	case ProdColThreshold2:
		if p.Value6408 != nil {
			if p.Value6408.Threshold2 {
				style.Image = assets.ImgPinOn
			} else {
				style.Image = assets.ImgPinOff
			}
		}
	case ProdColConnection:

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

func (m *DafProductsTable) Checked(index int) bool {
	return m.items[index].Checked
}

func (m *DafProductsTable) SetChecked(index int, checked bool) error {
	m.items[index].Checked = checked
	_, err := data.DBProducts.Exec(`UPDATE product SET checked = ? WHERE product_id = ?`, checked, m.items[index].ProductID)
	return err
}

type ProductColumn int

const (
	ProdColAddr ProductColumn = iota
	ProdColSerialNumber
	ProdColProductID
	ProdColConcentration
	ProdColCurrent
	ProdColMode
	ProdColFailure
	ProdColThreshold1
	ProdColThreshold2
	ProdColVersion
	ProdColGas
	ProdColConnection
)

var ProductColumns = func() []TableViewColumn {
	x := make([]TableViewColumn, ProdColConnection+1)

	type t = TableViewColumn
	x[ProdColAddr] =
		t{Title: "Адрес", Width: 80}
	x[ProdColSerialNumber] =
		t{Title: "Номер", Width: 80}
	x[ProdColProductID] =
		t{Title: "ID", Width: 80}
	x[ProdColConcentration] =
		t{Title: "Концентрация", Width: 150, Precision: 3}
	x[ProdColCurrent] =
		t{Title: "Ток", Width: 100, Precision: 1}
	x[ProdColThreshold1] =
		t{Title: "Порог 1", Width: 120, Precision: 1}
	x[ProdColThreshold2] =
		t{Title: "Порог 2", Width: 120, Precision: 1}
	x[ProdColMode] =
		t{Title: "Режим"}
	x[ProdColFailure] =
		t{Title: "Отказ"}
	x[ProdColVersion] =
		t{Title: "Версия"}
	x[ProdColGas] =
		t{Title: "Газ"}
	x[ProdColConnection] =
		t{Title: "Связь"}

	return x
}()