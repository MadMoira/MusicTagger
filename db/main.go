package db

import (
	"log"

	"github.com/jmoiron/sqlx"

	// Init the driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE test (
	id integer,
	name text,
	path text,
	primary key (id)
)
`

// Db main db usage
var Db *sqlx.DB

// InitDb created the initial schema
func InitDb() {
	db, err := sqlx.Connect("sqlite3", "./musictagger.db")
	if err != nil {
		log.Panic("Error starting the database")
	}

	db.MustExec(schema)
}
