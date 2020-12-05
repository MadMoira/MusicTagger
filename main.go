package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/mikkyang/id3-go"
	"github.com/rivo/tview"

	"os"
)

var app *tview.Application
var list *tview.List
var tw *tview.TextView
var frm *tview.Form

var files []string
var infos []os.FileInfo
var currentPath string

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	app = tview.NewApplication()
	list = tview.NewList()
	list.SetBorder(true)
	list.SetInputCapture(listEventHandler)
	tw = tview.NewTextView()
	tw.SetBorder(true)
	tw.SetText("Hello World")
	frm = tview.NewForm()

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewFlex().
					AddItem(list, 0, 1, true).
					AddItem(tw, 0, 1, false), 0, 7, true).
				AddItem(tview.NewTextView().SetText("testing"), 0, 1, false), 0, 1, false)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	currentPath = path

	// err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
	// 	files = append(files, path)
	// 	infos = append(infos, info)
	// 	return nil
	// })

	retrieveDirFiles(path)
	addPathsToList(files)

	if err := app.SetRoot(flex, true).SetInputCapture(appEventHandler).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func retrieveDirFiles(path string) {
	files = nil
	infos = nil

	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirFiles {
		files = append(files, path+"/"+file.Name())
		infos = append(infos, file)
	}
}

func addPathsToList(files []string) {
	for _, file := range files {
		list.AddItem(file, "", 0, nil)
	}
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

func showMetadata(path string) {

	m, err := id3.Open(path)
	if err != nil {
		log.Fatalf("Error loading file %v", path)
		return
	}
	defer m.Close()

	// log.Printf("Album %v", m.Album())
	// log.Printf("Artis %v", m.Artist())
	log.Printf("TALB %v", m.Frame("TALB"))
	log.Printf("TCOM %v", m.Frame("TCOM"))
	log.Printf("TEXT %v", m.Frame("TEXT"))
	log.Printf("TPE2 %v", m.Frame("TPE2"))
	log.Printf("Artist %v", m.AllFrames())

}

func listEventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Key() == tcell.KeyEnter {
		currentItem := list.GetCurrentItem()
		if infos[currentItem].IsDir() {
			log.Printf("Checking path %v", files[currentItem])
			currentPath = files[currentItem]
			retrieveDirFiles(currentPath)
			list.Clear()
			addPathsToList(files)
			tw.SetText("New Dir")
		} else {
			showMetadata(files[currentItem])
			tw.SetTextColor(tcell.ColorBlue)
			tw.SetText(strconv.FormatBool(infos[currentItem].IsDir()))
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
	}

	return eventKey
}
