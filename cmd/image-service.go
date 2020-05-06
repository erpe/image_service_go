package main

import (
	"github.com/erpe/image_service_go/app"
	"github.com/erpe/image_service_go/app/config"
	"log"
)

func main() {
	log.Println("Hello ImageService!")
	log.Println("Uses Bearer Token: ", config.GetConfig().Server.Token)
	app := app.App{}
	app.Initialize(config.GetConfig())
	app.Run()
}
