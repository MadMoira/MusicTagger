package gui

import (
	"bytes"
	"log"
	"musictagger/core"
	"text/template"

	"github.com/rivo/tview"

	"os"
)

var app *tview.Application
var bodyFlex *tview.Flex

var list *tview.List
var detailTw *tview.TextView
var topTw *tview.TextView
var bottomTw *tview.TextView

var frm *tview.Form

var files []string
var infos []os.FileInfo
var currentPath string

var oldSelection int

var helpText string = `[S] Store current metadata	[R] Recover metadata	[s] Store current song
[alt+s] Store all folder songs (including subfolders)	[alt+s] Quit app`

// Init Starts the GUI
func Init() {
	list = tview.NewList()
	list.SetBorder(true)
	list.SetInputCapture(listEventHandler)
	detailTw = tview.NewTextView()
	detailTw.SetBorder(true)
	detailTw.SetText("")
	topTw = tview.NewTextView()
	bottomTw = tview.NewTextView()
	bottomTw.SetText(helpText)

	app = tview.NewApplication()
	bodyFlex = tview.NewFlex()
	bodyFlex.AddItem(list, 0, 2, false)
	bodyFlex.AddItem(detailTw, 0, 3, false)

	frm = tview.NewForm()
	frm.AddInputField("TALB", "", 30, nil, nil)
	frm.AddInputField("TIT2", "", 30, nil, nil)
	frm.AddInputField("TPE1", "", 30, nil, nil)
	frm.AddInputField("TPE2", "", 30, nil, nil)
	frm.AddInputField("TCON", "", 30, nil, nil)
	frm.AddInputField("TRCK", "", 30, nil, nil)
	frm.AddInputField("TYER", "", 30, nil, nil)
	frm.AddButton("Save", func() {
		app.Stop()
	})

	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(topTw, 0, 1, false).
				AddItem(bodyFlex, 0, 10, false).
				AddItem(bottomTw, 0, 1, false), 0, 1, false)

	currentPath = core.Settings.Path

	retrieveDirFiles(currentPath)
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

func showMetadata(songName string) {

	log.Printf("Reading metadata for: %v", songName)

	song := getSongMetadata(songName)

	tpl, err := template.ParseFiles("gui/data.tpl")
	if err != nil {
		log.Panic("Failed to retrieve template")
		return
	}

	var result bytes.Buffer
	if err := tpl.Execute(&result, song); err != nil {
		log.Panicf("Error while rendering tempalte %v", err)
	}

	detailTw.SetText(result.String())
}
