package main

import (
	"log"

	"github.com/pocketbase/pocketbase"

	"videoeditor/src"
	"videoeditor/src/config"
	"videoeditor/src/modelhooks"
)

func main() {
	config.LoadConfigs()

	app := pocketbase.New()

	src.SetupRouter(app)

	modelhooks.Register(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
