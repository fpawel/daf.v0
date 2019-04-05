package data

import (
	"database/sql"
	"github.com/fpawel/elco/pkg/winapp"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/sqlite3"
	"path/filepath"
)

//go:generate go run github.com/fpawel/elco/cmd/utils/sqlstr/...

type Gas int

const (
	Gas1 Gas = 1
	Gas2 Gas = 2
	Gas3 Gas = 3
	Gas4 Gas = 4
)

func GetLastParty() (party *Party) {
	party = new(Party)
	err := DBProducts.SelectOneTo(party, `ORDER BY created_at DESC LIMIT 1;`)
	if err == reform.ErrNoRows {
		partyID := CreateNewParty()
		err = DBProducts.FindByPrimaryKeyTo(party, partyID)
	}
	if err != nil {
		panic(err)
	}
	return
}

func GetLastPartyID() (partyID int64) {
	row := DBProducts.QueryRow(`SELECT party_id FROM party ORDER BY created_at DESC LIMIT 1`)
	if err := row.Scan(&partyID); err == sql.ErrNoRows {
		CreateNewParty()
	}
	return
}

func GetProductsOfLastParty() (products []*Product) {
	xs, err := DBProducts.SelectAllFrom(
		ProductTable, "WHERE party_id = ? ORDER BY created_at", GetLastPartyID())
	if err != nil {
		panic(err)
	}
	for _, x := range xs {
		p := x.(*Product)
		products = append(products, p)
	}
	return
}

func GetProductsByPartyID(partyID int64) (products []*Product) {
	xs, err := DBProducts.SelectAllFrom(
		ProductTable,
		"WHERE party_id = ? ORDER BY created_at", partyID)
	if err != nil {
		panic(err)
	}
	for _, x := range xs {
		p := x.(*Product)
		products = append(products, p)
	}
	return
}

func HasCheckedProducts(products []*Product) bool {
	for _, p := range products {
		if p.Checked {
			return true
		}
	}
	return false
}

func CreateNewParty() int64 {
	r, err := DBProducts.Exec(`INSERT INTO party DEFAULT VALUES`)
	if err != nil {
		panic(err)
	}
	partyID, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	return partyID
}

var (
	//DBxProducts *sqlx.DB
	DBProducts *reform.DB
)

func init() {

	dataFolder, err := winapp.AppDataFolderPath()
	if err != nil {
		panic(err)
	}
	dataFolder = filepath.Join(dataFolder, "daf")
	err = winapp.EnsuredDirectory(dataFolder)
	if err != nil {
		panic(err)
	}
	fileName := filepath.Join(dataFolder, "daf.sqlite")

	conn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		panic(err)
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)
	conn.SetConnMaxLifetime(0)

	if _, err = conn.Exec(SQLCreate); err != nil {
		panic(err)
	}

	//DBxProducts = sqlx.NewDb(conn, "sqlite3")
	DBProducts = reform.NewDB(conn, sqlite3.Dialect, nil)
}
