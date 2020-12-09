package gui

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var validExtensions []string

func init() {
	validExtensions = []string{"mp3"}
}

func retrieveAllFiles(path string) {
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
		log.Fatalf("Error retrieving all files")
	}

}

func retrieveDirFiles(path string) {
	files = nil
	infos = nil

	topTw.SetText(path)
	dirFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
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
