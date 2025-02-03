package repo

import (
	"bearguard/cm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path"
)

func getDbMust() *sqlx.DB {
	dbPath := path.Join(cm.GetProjectRoot(), "liveguard.db")
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Panic(err)
	}
	return db
}
