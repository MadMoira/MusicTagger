package main

import (
	"flag"
	"log"
	"os"

	"musictagger/db"
	"musictagger/gui"
)

func main() {
	f := setMainLogger()
	defer f.Close()

	initDb := flag.Bool("db", false, "Set if you require to start the local db")
	removeDb := flag.Bool("rmdb", false, "Remove the localdb")
	flag.Parse()

	if *removeDb {
		err := os.Remove("./musictagger.db")
		if err != nil {
			log.Print("Database didn't exist")
		}
	}

	if *initDb {
		db.InitDb()
	}

	if !*removeDb && !*initDb {
		gui.Init()
	}
}

func setMainLogger() *os.File {
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
