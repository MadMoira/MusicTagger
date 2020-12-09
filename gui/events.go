package gui

import (
	"log"
	"strings"

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
			showMetadata(files[currentItem])
		}
	}

	if eventKey.Rune() == 'a' {
		a := strings.Split(currentPath, "/")
		a = a[:len(a)-1]
		backPath := strings.Join(a, "/")
		retrieveDirFiles(backPath)
		currentPath = backPath
		list.Clear()
		addPathsToList(files)
		list.SetCurrentItem(oldSelection)
	}

	return eventKey
}

func appEventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Rune() == 'q' {
		app.Stop()
	}

	if eventKey.Rune() == 'e' {
		app.SetFocus(list)
	}
	return eventKey
}
