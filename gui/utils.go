package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"musictagger/core"
	"musictagger/db"

	"github.com/mikkyang/id3-go"

	v2 "github.com/mikkyang/id3-go/v2"
)

var validExtensions []string

func init() {
	validExtensions = []string{"mp3"}
}

func retrieveAllFiles(path string) ([]string, []os.FileInfo) {
	var files []string
	var infos []os.FileInfo

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		pathParts := strings.Split(info.Name(), ".")
		ext := pathParts[len(pathParts)-1]

		found := false
		for _, validExt := range validExtensions {
			if ext == validExt {
				found = true
				break
			}
		}

		if found {
			files = append(files, path)
			infos = append(infos, info)
		}

		return nil
	})

	if err != nil {
		log.Printf("Error retrieving all files")
		return nil, nil
	}

	return files, infos
}

func retrieveDirFiles(path string) {
	files = nil
	infos = nil

	topTw.SetText(path)
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	for _, file := range dirFiles {
		pathParts := strings.Split(file.Name(), ".")
		ext := pathParts[len(pathParts)-1]

		found := false
		for _, validExt := range validExtensions {
			if ext == validExt {
				found = true
				break
			}
		}

		if found || file.IsDir() {
			files = append(files, file.Name())
			infos = append(infos, file)
		}
	}
}

func getSongMetadata(songPath string) *core.SongMetadata {

	m, err := id3.Open(songPath)
	if err != nil {
		log.Printf("Error loading file %v", songPath)
		return nil
	}
	defer m.Close()

	log.Printf("Reading metadata for: %v", songPath)

	song := core.SongMetadata{
		Path: songPath,
		TALB: fmt.Sprintf("%v", m.Frame("TALB")),
		TIT2: fmt.Sprintf("%v", m.Frame("TIT2")),
		TPE1: fmt.Sprintf("%v", m.Frame("TPE1")),
		TPE2: fmt.Sprintf("%v", m.Frame("TPE2")),
		TCON: fmt.Sprintf("%v", m.Frame("TCON")),
		TRCK: fmt.Sprintf("%v", m.Frame("TRCK")),
		TYER: fmt.Sprintf("%v", m.Frame("TYER")),
	}

	return &song
}

func editFrame(f *id3.File, frame string, newValue string) {
	if len(newValue) == 0 {
		return
	}

	f.DeleteFrames(frame)
	createFrame(f, frame, newValue)
}

func createFrame(f *id3.File, frame string, newValue string) {
	if len(newValue) == 0 {
		return
	}

	ft := v2.V23FrameTypeMap[frame]
	textFrame := v2.NewTextFrame(ft, "")
	textFrame.SetEncoding("UTF-8")
	textFrame.SetText(newValue)
	f.AddFrames(textFrame)
}

func saveSong(newMetadata core.SongMetadata) {
	log.Printf("Editing file %v", newMetadata)

	songPath := newMetadata.Path
	m, err := id3.Open(songPath)
	if err != nil {
		log.Printf("Error loading file %v", songPath)
		return
	}
	defer m.Close()

	frameList := []string{"TALB", "TIT2", "TPE1", "TPE2", "TCON", "TRCK", "TYER"}
	var newValues map[string]string

	in, _ := json.Marshal(newMetadata)
	json.Unmarshal(in, &newValues)

	for _, frame := range frameList {
		if m.Frame(frame) != nil {
			if frame == "TRCK" || frame == "TYER" {
				editFrame(m, frame, newValues[frame])
			} else {
				editFrame(m, frame, newValues[frame])
			}
		} else {
			createFrame(m, frame, newValues[frame])
		}
	}
}

func recoverSingleSong(songPath string) {
	log.Printf("Recovering file %v", songPath)

	m, err := id3.Open(songPath)
	if err != nil {
		log.Printf("Error loading file %v", songPath)
		return
	}
	defer m.Close()

	sm := db.RecoverByPath(songPath)

	frameList := []string{"TALB", "TIT2", "TPE1", "TPE2", "TCON", "TRCK", "TYER"}
	var newValues map[string]string

	in, _ := json.Marshal(&sm)
	json.Unmarshal(in, &newValues)

	for _, frame := range frameList {
		if m.Frame(frame) != nil {
			if frame == "TRCK" || frame == "TYER" {
				editFrame(m, frame, newValues[frame])
			} else {
				editFrame(m, frame, newValues[frame])
			}
		} else {
			createFrame(m, frame, newValues[frame])
		}
	}
}

func debugFrames(frames []v2.Framer) {
	for _, frame := range frames {
		log.Printf("Frame Name %v", frame.Id())
		log.Printf("Frame Value %v", frame.String())
	}
}
