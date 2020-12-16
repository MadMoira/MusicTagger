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

var helpText string = `[s] Store song metadata [S] Store folder metadata [r] Recover song metadata [R] Recover folder metadata [e] Open single song edit [w] Focus list [Q] Quit app`

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
	frm.AddInputField("TALB", "", 50, nil, nil)
	frm.AddInputField("TIT2", "", 50, nil, nil)
	frm.AddInputField("TPE1", "", 50, nil, nil)
	frm.AddInputField("TPE2", "", 50, nil, nil)
	frm.AddInputField("TCON", "", 50, nil, nil)
	frm.AddInputField("TRCK", "", 50, nil, nil)
	frm.AddInputField("TYER", "", 50, nil, nil)
	frm.AddButton("Save", nil)

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

func formEditSingleSong(songPath string) {

	song := getSongMetadata(songPath)

	frm.GetFormItemByLabel("TALB").(*tview.InputField).SetText(song.TALB)
	frm.GetFormItemByLabel("TCON").(*tview.InputField).SetText(song.TCON)
	frm.GetFormItemByLabel("TIT2").(*tview.InputField).SetText(song.TIT2)
	frm.GetFormItemByLabel("TPE1").(*tview.InputField).SetText(song.TPE1)
	frm.GetFormItemByLabel("TPE2").(*tview.InputField).SetText(song.TPE2)
	frm.GetFormItemByLabel("TRCK").(*tview.InputField).SetText(song.TRCK)
	frm.GetFormItemByLabel("TYER").(*tview.InputField).SetText(song.TYER)

	frm.GetButton(0).SetSelectedFunc(func() {
		metadata := core.SongMetadata{
			Path: song.Path,
			TALB: frm.GetFormItemByLabel("TALB").(*tview.InputField).GetText(),
			TCON: frm.GetFormItemByLabel("TCON").(*tview.InputField).GetText(),
			TIT2: frm.GetFormItemByLabel("TIT2").(*tview.InputField).GetText(),
			TPE1: frm.GetFormItemByLabel("TPE1").(*tview.InputField).GetText(),
			TPE2: frm.GetFormItemByLabel("TPE2").(*tview.InputField).GetText(),
			TRCK: frm.GetFormItemByLabel("TRCK").(*tview.InputField).GetText(),
			TYER: frm.GetFormItemByLabel("TYER").(*tview.InputField).GetText(),
		}
		log.Print(metadata)
		editSingleSong(metadata)
	})

	bodyFlex.RemoveItem(detailTw)
	bodyFlex.RemoveItem(frm)
	bodyFlex.AddItem(frm, 0, 3, false)

	app.SetFocus(frm)
}

func showMetadata(songPath string) {

	song := getSongMetadata(songPath)

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

	bodyFlex.RemoveItem(detailTw)
	bodyFlex.RemoveItem(frm)
	bodyFlex.AddItem(detailTw, 0, 3, false)
}
