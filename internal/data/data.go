package data

import (
	"database/sql"
	"github.com/fpawel/elco/pkg/winapp"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

//go:generate go run github.com/fpawel/elco/cmd/utils/sqlstr/...

func LastParty() (party *Party) {
	party = new(Party)
	err := DBProducts.Get(party, `SELECT * FROM last_party`)
	if err == nil {
		return
	}
	if err != sql.ErrNoRows {
		panic(err)
	}
	CreateNewParty()
	if err = DBProducts.Get(&party, `SELECT * FROM last_party`); err != nil {
		panic(err)
	}
	return
}

func GetProductsByPartyID(partyID int64) (products []*Product) {
	if err := DBProducts.Select(&products, `SELECT * FROM product WHERE party_id = ?`, partyID); err != nil {
		panic(err)
	}
	return
}

func GetProductsOfLastParty() (products []*Product) {
	if err := DBProducts.Select(&products, `SELECT * FROM product WHERE party_id = (SELECT party_id FROM last_party)`); err != nil {
		panic(err)
	}
	return
}

func GetProductByID(productID int64) (product *Product, err error) {
	product = new(Product)
	err = DBProducts.Get(product, `SELECT * FROM product WHERE product_id = ?`, productID)
	return
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

var DBProducts *sqlx.DB

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

	DBProducts = sqlx.NewDb(conn, "sqlite3")
}
