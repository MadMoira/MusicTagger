package main

import (
	"flag"

	"musictagger/db"
	"musictagger/gui"
)

func main() {

	initDb := flag.Bool("db", false, "Set if you require to start the local db")
	flag.Parse()

	if *initDb {
		db.InitDb()
	} else {
		gui.Init()
	}
}
