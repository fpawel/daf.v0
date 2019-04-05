package viewmodel

import (
	. "github.com/lxn/walk/declarative"
)

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
