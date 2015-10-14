package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/compotlab/synopsis/app"
	"github.com/compotlab/synopsis/src"
	"io/ioutil"
	"net/http"
	"strings"
)

func RegisterRepoController(router *mux.Router) {
	router.Methods("GET").Path("/repo").HandlerFunc(RepoAllHandler)
	router.Methods("POST").Path("/repo").HandlerFunc(RepoSaveHandler)
	router.Methods("DELETE").Path("/repo").HandlerFunc(RepoDeleteHandler)
}

func RepoAllHandler(res http.ResponseWriter, req *http.Request) {
	config := prepareConfig()
	j, err := json.Marshal(config.File.Repositories)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(j)
}

func RepoSaveHandler(res http.ResponseWriter, req *http.Request) {
	config := prepareConfig()
	repo := src.Repository{}
	body, _ := ioutil.ReadAll(req.Body)
	if err := json.Unmarshal(body, &repo); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	urlRepo := strings.TrimSpace(repo.Url)
	typeRepo := strings.TrimSpace(repo.Type)
	originalUrlRepo := req.PostForm.Get("original_url")
	isExist := false
	for key, value := range config.File.Repositories {
		if value.Url == originalUrlRepo {
			isExist = true
			config.File.Repositories[key].Url = urlRepo
			config.File.Repositories[key].Type = typeRepo
		}
	}
	if !isExist {
		config.File.Repositories = append(config.File.Repositories, src.Repository{Type: typeRepo, Url: urlRepo})
	}
	err := saveNewConfig(config)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RepoDeleteHandler(res http.ResponseWriter, req *http.Request) {
	config := prepareConfig()
	url := req.URL.Query().Get("url")
	isExist := false
	for key, value := range config.File.Repositories {
		if value.Url == url {
			isExist = true
			config.File.Repositories = append(config.File.Repositories[:key], config.File.Repositories[key+1:]...)
		}
	}
	if !isExist {
		http.Error(res, "Item not exist!", http.StatusNoContent)
		return
	}
	err := saveNewConfig(config)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveNewConfig(config src.Config) error {
	j, err := json.MarshalIndent(config.File, "", "  ")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(config.FileName, j, 0755); err != nil {
		return err
	}
	return nil
}

func prepareConfig() src.Config {
	app := app.GetApp()
	config := src.Config{}
	config.PrepareConfig(app)
	return config
}
