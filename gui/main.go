package gui

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/mikkyang/id3-go"
	v2 "github.com/mikkyang/id3-go/v2"
	"github.com/rivo/tview"

	"os"
)

var app *tview.Application
var list *tview.List
var tw *tview.TextView
var topTw *tview.TextView
var bottomTw *tview.TextView
var frm *tview.Form

var files []string
var infos []os.FileInfo
var currentPath string

var oldSelection int

// Init Starts the GUI
func Init() {
	app = tview.NewApplication()
	list = tview.NewList()
	list.SetBorder(true)
	list.SetInputCapture(listEventHandler)
	tw = tview.NewTextView()
	tw.SetBorder(true)
	tw.SetText("")
	topTw = tview.NewTextView()
	bottomTw = tview.NewTextView()
	bottomTw.SetText("[S] Store current metadata\t[R] Recover metadata")

	frm = tview.NewForm()

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(topTw, 0, 1, false).
				AddItem(
					tview.NewFlex().
						AddItem(list, 0, 2, false).
						AddItem(tw, 0, 3, false),
					0, 10, false).
				AddItem(bottomTw, 0, 1, false), 0, 1, false)

	// path := "/mnt/f/Music/Afterglow"
	path := "/home/camilo.r/Documents/pcode/go/MusicTagger/testdata/Afterglow"
	currentPath = path

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

func addPathsToList(files []string) {
	list.Clear()
	for _, file := range files {
		list.AddItem(file, "", 0, nil)
	}
}

func debugFrames(frames []v2.Framer) {
	for _, frame := range frames {
		log.Printf("Frame Name %v", frame.Id())
		log.Printf("Frame Value %v", frame.String())
	}
}

func showMetadata(songName string) {

	m, err := id3.Open(currentPath + "/" + songName)
	if err != nil {
		log.Fatalf("Error loading file %v", songName)
		return
	}
	defer m.Close()

	log.Printf("Reading metadata for: %v", songName)

	tpl, err := template.ParseFiles("gui/data.tpl")
	if err != nil {
		log.Fatal("Failed to retrieve template")
		return
	}

	type metadata struct {
		TALB string
		TIT2 string
		TPE1 string
		TPE2 string
		TCON string
		TRCK string
		TYER string
	}

	song := metadata{
		TALB: fmt.Sprintf("%v", m.Frame("TALB")),
		TIT2: fmt.Sprintf("%v", m.Frame("TIT2")),
		TPE1: fmt.Sprintf("%v", m.Frame("TPE1")),
		TPE2: fmt.Sprintf("%v", m.Frame("TPE2")),
		TCON: fmt.Sprintf("%v", m.Frame("TCON")),
		TRCK: fmt.Sprintf("%v", m.Frame("TRCK")),
		TYER: fmt.Sprintf("%v", m.Frame("TYER")),
	}

	log.Printf("things %v", song)

	var result bytes.Buffer

	if err := tpl.Execute(&result, song); err != nil {
		panic(err)
	}

	tw.SetText(result.String())
}
