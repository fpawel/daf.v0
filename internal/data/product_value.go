package data

import (
	"time"
)

//go:generate reform

// ProductValue represents a row in product_value table.
//reform:product_value
type ProductValue struct {
	ProductValueID int64     `reform:"product_value_id,pk"`
	ProductID      int64     `reform:"product_id"`
	CreatedAt      time.Time `reform:"created_at"`
	Gas            Gas       `reform:"gas"`
	Name           string    `reform:"name"`
	Concentration  float64   `reform:"concentration"`
	Current        float64   `reform:"current"`
	Threshold1     bool      `reform:"threshold1"`
	Threshold2     bool      `reform:"threshold2"`
}
