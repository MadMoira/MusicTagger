package db

import (
	"log"

	"musictagger/core"

	"github.com/jmoiron/sqlx"

	// Init the driver for sqlite
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS music_info (
	id integer,
	path text,
	TALB text,
	TIT2 text,
	TPE1 text,
	TPE2 text,
	TCON text,
	TRCK text,
	TYER text,
	PRIMARY KEY (id),
	UNIQUE(path)
)
`

var dbConnection *sqlx.DB

// Connect creates a new DB Connection
func Connect() {
	var err error
	dbConnection, err = sqlx.Connect("sqlite3", "./musictagger.db")
	if err != nil {
		log.Panic("Error starting the database")
	}
	dbConnection.Ping()
	log.Print("Creating a new DB Connection")
}

// CloseSession close connection
func CloseSession() {
	dbConnection.Close()
}

// InitDb created the initial schema
func InitDb() {
	dbConnection.MustExec(schema)
}

// StoreSongData Stores all data sent into the db
func StoreSongData(sm core.SongMetadata) {
	storeSongStmt := `INSERT INTO music_info (path, TALB, TIT2, TPE1, TPE2, TCON, TRCK, TYER) 
					  VALUES (:path, :talb, :tit2, :tpe1, :tpe2, :tcon, :trck, :tyer)
					   ON CONFLICT(path) DO NOTHING`

	_, err := dbConnection.NamedExec(storeSongStmt, sm)
	if err != nil {
		log.Printf("Error while parsing storing the data %v", err)
		return
	}

	log.Printf("Storing data for %v", sm.Path)
}
