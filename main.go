package main

import (
	"github.com/compotlab/synopsis/app"
	"github.com/compotlab/synopsis/app/web"
)

func main() {
	app.NewApp()
	web.NewServer()
}
