package app

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Application struct {
	Config Config
}

type Config struct {
	Server map[string]string `yaml:"server"`
	Build  map[string]string `yaml:"build"`
}

var App Application

func NewApp() Application {
	filename, _ := filepath.Abs("./app.yaml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
	App = Application{Config: config}
	return App
}

func GetApp() Application {
	return App
}
