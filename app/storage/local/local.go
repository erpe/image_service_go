package local

import (
	"github.com/erpe/image_service_go/app/config"
	"image"
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

func ReadImage(fname string) (image.Image, string, error) {

	var img image.Image
	var format string

	filepath := DIRECTORY + "/" + fname

	infile, err := os.Open(filepath)

	if err != nil {
		return img, format, err
	}

	defer infile.Close()

	img, format, decErr := image.Decode(infile)

	if decErr != nil {
		return img, format, decErr
	}

	return img, format, nil
}

func ReadImageBytes(fname string) ([]byte, error) {

	var data []byte

	filepath := DIRECTORY + "/" + fname

	infile, err := os.Open(filepath)

	if err != nil {
		return data, err
	}

	defer infile.Close()

	data, readErr := ioutil.ReadAll(infile)

	if readErr != nil {
		return data, readErr
	}

	return data, nil
}
