package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/compotlab/synopsis/src"
	"github.com/compotlab/synopsis/src/packages"
	"io/ioutil"
	"net/http"
	"path"
)

var Lock bool

func RegisterPackageController(router *mux.Router) {
	router.HandleFunc("/package/all", AllPackagesHandler)
	router.HandleFunc("/package/update", PackageUpdateHandler)
}

func AllPackagesHandler(res http.ResponseWriter, req *http.Request) {
	config := &src.Config{}
	config.PrepareConfig()

	file, _ := ioutil.ReadFile(path.Join(config.OutputDir, "packages.json"))

	res.Header().Set("Content-Type", "application/json")
	res.Write(file)
}

func PackageUpdateHandler(res http.ResponseWriter, req *http.Request) {
	if !Lock {
		Lock = true

		config := &src.Config{}
		config.PrepareConfig()
		config.MakeOutputDir()

		flag := make(chan bool, config.ThreadNumber)
		done := make(chan int)

		pm := make(map[string]map[string]packages.Composer)
		i := 0
		for _, repo := range config.File.Repositories {
			go func(r src.Repository) {
				flag <- true
				r.Run(pm, config)
				defer func() {
					<-flag
					i += 1
					done <- i
				}()
			}(repo)
		}

		setEventStream(res, len(config.File.Repositories), done)

		p := packages.Packages{Package: pm}
		if err := p.ToJson(config.OutputDir); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		defer func() {
			Lock = false
			close(flag)
			close(done)
		}()
	}
}

func setEventStream(res http.ResponseWriter, l int, done chan int) {
	flusher, ok := res.(http.Flusher)
	if !ok {
		http.Error(res, "Streaming unsupported!", http.StatusInternalServerError)
	} else {
		res.Header().Set("Content-Type", "text/event-stream")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("Connection", "keep-alive")
		res.Header().Set("Access-Control-Allow-Origin", "*")
		for i := 0; i < l; i++ {
			id := <-done
			fmt.Fprintf(res, "id: %s \n", "update")
			fmt.Fprintf(res, "data: %d \n\n", id)
			flusher.Flush()
		}
	}
}
