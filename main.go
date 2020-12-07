package main

import (
	"fmt"
	"io/ioutil"
	"log"

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
					AddItem(list, 0, 1, false).
					AddItem(tw, 0, 1, false), 0, 7, false).
				AddItem(tview.NewTextView().SetText("testing"), 0, 1, false), 0, 1, false)

	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	path := "/mnt/f/Music/Afterglow"
	currentPath = path

	// err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
	// 	files = append(files, path)
	// 	infos = append(infos, info)
	// 	return nil
	// })

	retrieveDirFiles(path)
	addPathsToList(files)

	app.SetRoot(flex, true).
		SetFocus(list).
		SetInputCapture(appEventHandler).
		EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func retrieveDirFiles(path string) {
	files = nil
	infos = nil

	list.SetTitle(path)
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirFiles {
		files = append(files, file.Name())
		infos = append(infos, file)
	}
}

func addPathsToList(files []string) {
	list.Clear()
	for _, file := range files {
		list.AddItem(file, "", 0, nil)
	}
}

func showMetadata(songName string) {

	m, err := id3.Open(currentPath + "/" + songName)
	if err != nil {
		log.Fatalf("Error loading file %v", songName)
		return
	}
	defer m.Close()

	template := fmt.Sprintf(
		`
	TALB: %v
	TCOM: %v
	TEXT: %v
	TPE1: %v
	TPE2: %v
	`, m.Frame("TALB"), m.Frame("TCOM"), m.Frame("TEXT"), m.Frame("TPE1"), m.Frame("TPE2"))

	tw.SetText(template)

	// log.Printf("Album %v", m.Album())
	// log.Printf("Artis %v", m.Artist())

	frames := m.AllFrames()

	for _, frame := range frames {
		log.Printf("%v %v %v", frame.Id(), frame.Size(), frame.String())
	}

	// log.Printf("TALB %v", m.Frame("TALB"))
	// log.Printf("TCOM %v", m.Frame("TCOM"))
	// log.Printf("TEXT %v", m.Frame("TEXT"))
	// log.Printf("TPE2 %v", m.Frame("TPE2"))

}
