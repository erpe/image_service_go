package local

import (
	"github.com/erpe/image_service_go/app/config"
	"io/ioutil"
)

var DIRECTORY, HOST string

func init() {
	HOST = config.GetConfig().Localstore.Assethost
	DIRECTORY = config.GetConfig().Localstore.Directory
}

func StoreLocal(buffer []byte, fname string) (string, error) {

	err := ioutil.WriteFile(DIRECTORY+"/"+fname, buffer, 0644)

	if err != nil {
		return "", err
	}

	url := HOST + "/" + DIRECTORY + "/" + fname
	return url, nil
}
