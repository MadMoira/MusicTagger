package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	// Init the driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE music_info (
	id integer,
	path text,
	TALB text,
	TIT2 text,
	TPE1 text,
	TPE2 text,
	TCON text,
	TRCK text,
	TYER text,
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

// StoreFileInfoData Stores all data sent by the 
func StoreFileInfoData(infos []os.FileInfo) {

}
