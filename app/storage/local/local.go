package local

import (
	"github.com/erpe/image_service_go/app/config"
	"io/ioutil"
	"os"
)

var DIRECTORY, HOST string

func init() {
	HOST = config.GetConfig().Localstore.Assethost
	DIRECTORY = config.GetConfig().Localstore.Directory
}

func SaveImage(buffer []byte, fname string) (string, error) {

	err := ioutil.WriteFile(DIRECTORY+"/"+fname, buffer, 0644)

	if err != nil {
		return "", err
	}

	url := HOST + "/" + DIRECTORY + "/" + fname
	return url, nil
}

func UnlinkImage(fname string) error {
	return os.Remove(DIRECTORY + "/" + fname)
}
