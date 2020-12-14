package core

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

// SongMetadata Core structure for storing and displaying song
type SongMetadata struct {
	Path string
	TALB string
	TIT2 string
	TPE1 string
	TPE2 string
	TCON string
	TRCK string
	TYER string
}

type configuration struct {
	Path string `toml:"string"`
}

// Settings contains global app settings
var Settings configuration

// LoadConfiguration Tries to load configuration or defaults values
func LoadConfiguration() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic("Error while getting current folder")
	}

	Settings = configuration{
		Path: dir,
	}

	confPath := dir + "/" + "conf.toml"
	tomlConfig, err := toml.LoadFile(confPath)
	if err != nil {
		log.Printf("Error reading configuration file: missing or wrong permissions %v", confPath)
		return
	}

	if expectedPath := tomlConfig.Get("path"); expectedPath != nil {
		Settings.Path = expectedPath.(string)
	}
}
