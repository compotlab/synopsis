package web

import (
	"github.com/fogcreek/profiler"
	"github.com/gorilla/mux"
	"github.com/compotlab/synopsis/app"
	"github.com/compotlab/synopsis/app/controller"
	"log"
	"net/http"
	"net/http/pprof"
	"path"
)

func NewServer() {
	app := app.GetApp()
	router := mux.NewRouter().StrictSlash(true)
	controller.RegisterRepoController(router)
	controller.RegisterPackageController(router)
	RegisterProfilerController(router)
	RegisterFileServer(app, router)
	if err := http.ListenAndServe(app.Config.Server["host"]+":"+app.Config.Server["port"], router); err != nil {
		log.Fatal(err)
	}
}

func RegisterFileServer(app app.Application, router *mux.Router) {
	outputDir := "./" + app.Config.Build["output"] + "/"
	router.HandleFunc("/packages.json", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, path.Join(outputDir, "packages.json"))
	})
	router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir(path.Join(outputDir, "dist")))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	router.PathPrefix("/admin").Handler(http.StripPrefix("/admin", http.FileServer(http.Dir("./app/view/admin/"))))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./app/view/package/"))))
}

func RegisterProfilerController(router *mux.Router) {
	router.HandleFunc("/profiler/info.html", profiler.MemStatsHTMLHandler)
	router.HandleFunc("/profiler/info", profiler.ProfilingInfoJSONHandler)
	router.HandleFunc("/profiler/start", profiler.StartProfilingHandler)
	router.HandleFunc("/profiler/stop", profiler.StopProfilingHandler)

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
}
