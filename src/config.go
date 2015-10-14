package src

import (
	"encoding/json"
	"github.com/compotlab/synopsis/app"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

type Config struct {
	File         File
	FileName     string
	ThreadNumber int
	OutputDir    string
	DistDir      string
}

type File struct {
	Name         string       `json:"name"`
	Homepage     string       `json:"homepage"`
	Archive      Archive      `json:"archive"`
	Repositories []Repository `json:"repositories"`
}

type Archive struct {
	Dir     string `json:"directory"`
	SkipDev bool   `json:"skip-dev"`
}

type Repository struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

func (config *Config) PrepareConfig(app app.Application) {
	config.FileName = app.Config.Build["file"]
	file, err := ioutil.ReadFile(config.FileName)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(file, &config.File); err != nil {
		log.Fatal(err)
	}
	config.ThreadNumber, _ = strconv.Atoi(app.Config.Build["thread"])
	config.OutputDir = app.Config.Build["output"]
	config.DistDir = path.Join(config.OutputDir, config.File.Archive.Dir)
}

func (config *Config) MakeOutputDir() {
	if _, err := os.Stat(config.OutputDir); err != nil {
		os.MkdirAll(config.OutputDir, 0777)
	}
}
