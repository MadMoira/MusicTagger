package gui

import (
	"bytes"
	"io/ioutil"
	"log"
	"text/template"

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

// Init Starts the GUI
func Init() {
	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
					AddItem(list, 0, 2, false).
					AddItem(tw, 0, 3, false), 0, 7, false).
				AddItem(tview.NewTextView().SetText("testing"), 0, 1, false), 0, 1, false)

	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	// path := "/mnt/f/Music/Afterglow"
	path := "/home/camilo.r/Documents/pcode/go/MusicTagger/testdata/Afterglow"
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
		TALB: m.Frame("TALB").String(),
		TIT2: m.Frame("TIT2").String(),
		TPE1: m.Frame("TPE1").String(),
		TPE2: m.Frame("TPE2").String(),
		TCON: m.Frame("TCON").String(),
		TRCK: m.Frame("TRCK").String(),
		TYER: m.Frame("TYER").String(),
	}

	log.Printf("things %v", song)

	var result bytes.Buffer

	if err := tpl.Execute(&result, song); err != nil {
		panic(err)
	}

	tw.SetText(result.String())
}
