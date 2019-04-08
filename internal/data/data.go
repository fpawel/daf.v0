package data

import (
	"database/sql"
	"github.com/fpawel/elco/pkg/winapp"
	"github.com/jmoiron/sqlx"
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

func GetProductEntry(productID int64, workName string) (ent *ProductEntry) {
	ent = new(ProductEntry)
	err := DBProducts.SelectOneTo(ent, "WHERE product_id = ? AND name = ?", workName)
	if err != reform.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return
}

func ClearCurrentProductsWorkResult(workName string) {
	DBxProducts.MustExec(`
DELETE FROM product_entry
WHERE name = ? 
AND product_id IN ( 
	SELECT product_id 
	FROM product 
	WHERE party_id = (SELECT party_id FROM last_party))`, workName)
}

func SetProductWorkInfo(productID int64, workName, message string) {
	SetProductEntry(productID, workName, true, message)
}
func SetProductWorkError(productID int64, workName string, err error) {
	SetProductEntry(productID, workName, false, err.Error())
}

func SetProductEntry(productID int64, workName string, ok bool, message string) {
	DBxProducts.MustExec(`
INSERT OR REPLACE INTO product_entry (product_id, name, ok, message)
VALUES (?, ?, ?, ?)`, productID, workName, ok, message)
}

var (
	DBxProducts *sqlx.DB
	DBProducts  *reform.DB
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

	DBxProducts = sqlx.NewDb(conn, "sqlite3")
	DBProducts = reform.NewDB(conn, sqlite3.Dialect, nil)
}
