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

CREATE TABLE IF NOT EXISTS folder_status (
	id integer,
	path text,
	visited boolean,
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
					  VALUES (:path, :TALB, :TIT2, :TPE1, :TPE2, :TCON, :TRCK, :TYER)
					   ON CONFLICT(path) DO NOTHING`

	_, err := dbConnection.NamedExec(storeSongStmt, sm)
	if err != nil {
		log.Printf("Error while parsing storing the data %v", err)
		return
	}

	log.Printf("Storing data for %v", sm.Path)
}

// StoreFolderStatus placeholder comment
func StoreFolderStatus(folderPath string) {
	storeFolderStmt := `INSERT INTO folder_status (path, visited)
						VALUES (?, ?) ON CONFLICT (path) DO NOTHING`

	_, err := dbConnection.Exec(storeFolderStmt, folderPath, true)
	if err != nil {
		log.Printf("Error while storing folder status %v", folderPath)
	}

	log.Printf("Folder %v visited", folderPath)
}

// RecoverSongByPath recover from the database one song by path
func RecoverSongByPath(songPath string) *core.SongMetadata {
	retrieveSongByPath := `SELECT path, talb, tit2, tpe1, tpe2, tcon, trck, tyer
						   FROM music_info WHERE path = ?`

	sm := core.SongMetadata{}
	err := dbConnection.Get(&sm, retrieveSongByPath, songPath)
	if err != nil {
		log.Printf("Error querying: %v", err)
		return nil
	}

	return &sm
}
