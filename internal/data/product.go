package data

import "time"

type Product struct {
	ProductID int64     `db:"product_id"`
	PartyID   int64     `db:"party_id"`
	CreatedAt time.Time `db:"created_at"`
	Addr      int       `db:"addr"`
	Serial    int64     `db:"serial"`
	Checked   bool      `db:"checked"`
}

func (p *Product) Save() error {
	_, err := DBProducts.NamedExec(
		`UPDATE product SET checked = :checked, serial = :serial, addr=:addr WHERE product_id = :product_id`,
		map[string]interface{}{
			"product_id": p.ProductID,
			"checked":    p.Checked,
			"serial":     p.Serial,
			"addr":       p.Addr,
		})
	return err
}
