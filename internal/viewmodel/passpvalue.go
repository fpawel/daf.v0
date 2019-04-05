package viewmodel

import (
	"github.com/fpawel/daf/internal/data"
	"github.com/lxn/walk"
)

type PassportValue struct {
	walk.TableModelBase
	product *data.Product
	values  []*data.ProductValue
}
