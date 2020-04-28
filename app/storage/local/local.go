package local

import (
	"github.com/erpe/image_service_go/app/config"
	"io/ioutil"
)

var DIRECTORY string

func init() {
	DIRECTORY = config.GetConfig().Localstore.Directory
}

func StoreLocal(buffer []byte, fname string) (string, error) {

	err := ioutil.WriteFile(DIRECTORY+"/"+fname, buffer, 0644)

	if err != nil {
		return "", err
	}

	return fname, nil
}
