package main

import (
	"github.com/austauschkompass/image_service_go/app"
	"github.com/austauschkompass/image_service_go/app/config"
	"log"
)

func main() {
	log.Println("Hello ImageService!")
	app := app.App{}
	app.Initialize(config.GetConfig())
	app.Run()
}
