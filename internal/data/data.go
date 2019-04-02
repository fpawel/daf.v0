package data

import (
	"database/sql"
	"github.com/fpawel/elco/pkg/winapp"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

//go:generate go run github.com/fpawel/elco/cmd/utils/sqlstr/...




func LastParty() (party Party) {
	err := DBProducts.Get(&party, `SELECT * FROM last_party`)
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

func SaveParty(party Party) error {
	_, err := DBProducts.NamedExec(
		`UPDATE party SET type = :type, pgs1 = :pgs1, pgs2 = :pgs2, pgs3 = :pgs3, pgs4 = :pgs4 WHERE party_id = :party_id`,
		map[string]interface{}{
			"party_id": party.PartyID,
			"pgs1":     party.Pgs1,
			"pgs2":     party.Pgs2,
			"pgs3":     party.Pgs3,
			"pgs4":     party.Pgs4,
			"type":     party.Type,
		})
	return err
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
