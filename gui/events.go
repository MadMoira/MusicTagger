package gui

import (
	"log"
	"strings"

	"musictagger/db"

	"github.com/gdamore/tcell/v2"
)

func listEventHandler(eventKey *tcell.EventKey) *tcell.EventKey {

	if eventKey.Key() == tcell.KeyEnter {
		currentItem := list.GetCurrentItem()
		if infos[currentItem].IsDir() {
			log.Printf("Checking path %v", files[currentItem])
			currentPath = currentPath + "/" + files[currentItem]
			oldSelection = currentItem
			retrieveDirFiles(currentPath)
			addPathsToList(files)
		} else {
			songPath := currentPath + "/" + files[currentItem]
			showMetadata(songPath)
		}
	}

	if eventKey.Key() == tcell.KeyBackspace2 {
		detailTw.SetText("")
		a := strings.Split(currentPath, "/")
		a = a[:len(a)-1]
		backPath := strings.Join(a, "/")
		retrieveDirFiles(backPath)
		currentPath = backPath
		addPathsToList(files)
		list.SetCurrentItem(oldSelection)
	}

	if eventKey.Rune() == 'e' {
		currentItem := list.GetCurrentItem()

		if !infos[currentItem].IsDir() {
			songPath := currentPath + "/" + files[currentItem]
			formEditSingleSong(songPath)
		}
	}

	if eventKey.Rune() == 'r' {
		currentItem := list.GetCurrentItem()
		if !infos[currentItem].IsDir() {
			songPath := currentPath + "/" + files[currentItem]
			recoverSingleSong(songPath)
		}
	}

	if eventKey.Rune() == 'R' {
		allPaths, _ := retrieveAllFiles(currentPath)
		for _, path := range allPaths {
			recoverSingleSong(path)
		}
	}

	if eventKey.Rune() == 's' {
		currentItem := list.GetCurrentItem()
		songPath := currentPath + "/" + files[currentItem]
		song := getSongMetadata(songPath)

		db.StoreSongData(*song)
		log.Printf("Finished storing song %v", files[currentItem])
	}

	if eventKey.Rune() == 'S' {
		log.Print("Storing all songs")
		log.Print(currentPath)
		allPaths, _ := retrieveAllFiles(currentPath)
		for _, path := range allPaths {
			song := getSongMetadata(path)
			db.StoreSongData(*song)
			log.Printf("Finished storing song %v", path)
		}
	}

	return eventKey
}

func appEventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Rune() == 'w' {
		app.SetFocus(list)
	}
	if eventKey.Rune() == 'Q' {
		app.Stop()
	}
	return eventKey
}
