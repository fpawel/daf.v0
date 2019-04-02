package data

import "time"

type Party struct {
	PartyID   int64     `db:"party_id"`
	CreatedAt time.Time `db:"created_at"`
	Type      int       `db:"type"`
	Pgs1      float64   `db:"pgs1"`
	Pgs2      float64   `db:"pgs2"`
	Pgs3      float64   `db:"pgs3"`
	Pgs4      float64   `db:"pgs4"`
}

func (p *Party) Pgs(gas Gas) float64 {
	switch gas {
	case Gas1:
		return p.Pgs1
	case Gas2:
		return p.Pgs2
	case Gas3:
		return p.Pgs3
	case Gas4:
		return p.Pgs4
	default:
		panic("wrong gas")
	}
}

func (p *Party) SetPgs(gas Gas, value float64) {
	switch gas {
	case Gas1:
		p.Pgs1 = value
	case Gas2:
		p.Pgs2 = value
	case Gas3:
		p.Pgs3 = value
	case Gas4:
		p.Pgs4 = value
	default:
		panic("wrong gas")
	}
}

func (p *Party) Save() error {
	_, err := DBProducts.NamedExec(
		`UPDATE party SET type = :type, pgs1 = :pgs1, pgs2 = :pgs2, pgs3 = :pgs3, pgs4 = :pgs4 WHERE party_id = :party_id`,
		map[string]interface{}{
			"party_id": p.PartyID,
			"pgs1":     p.Pgs1,
			"pgs2":     p.Pgs2,
			"pgs3":     p.Pgs3,
			"pgs4":     p.Pgs4,
			"type":     p.Type,
		})
	return err
}

type Gas int

const (
	Gas1 Gas = 1
	Gas2 Gas = 2
	Gas3 Gas = 3
	Gas4 Gas = 4
)
