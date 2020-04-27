package handler

import (
	"github.com/erpe/image_service_go/app/config"
	"log"
)

var DIRECTORY string

func init() {
	DIRECTORY = config.GetConfig().Localstore.Directory
}

func storeLocal(buffer []byte, fname string) (string, error) {
	log.Printf("buffer: %s fname: %s", buffer, fname)
	return fname, nil

}
