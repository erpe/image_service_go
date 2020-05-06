package storage

import (
	"errors"
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/storage/local"
	"github.com/erpe/image_service_go/app/storage/s3store"
	"image"
	"log"
)

func SaveImage(data []byte, name string) (string, error) {
	cfg := config.GetConfig()

	if cfg.Storage.IsS3() {
		res, err := s3store.SaveImage(data, name)

		if err != nil {
			return "", err
		}

		return res, nil
	}

	if cfg.Storage.IsLocal() {
		res, err := local.SaveImage(data, name)

		if err != nil {
			return "", err
		}

		return res, nil
	}

	return "", errors.New("Storage type unknown")
}

func UnlinkImage(fname string) error {
	cfg := config.GetConfig()

	if cfg.Storage.IsS3() {
		return s3store.UnlinkImage(fname)
	}

	if cfg.Storage.IsLocal() {
		return local.UnlinkImage(fname)
	}

	return errors.New("Storage type unknown")
}

func ReadImage(fname string) (image.Image, string, error) {
	cfg := config.GetConfig()

	var img image.Image
	var format string

	if cfg.Storage.IsS3() {

		img, format, err := s3store.ReadImage(fname)

		if err != nil {
			log.Println("ERROR - storage.ReadImage: ", err.Error())
			return img, format, err
		}
		return img, format, nil
	}

	if cfg.Storage.IsLocal() {

		img, format, err := local.ReadImage(fname)

		if err != nil {
			return img, format, err
		}
		return img, format, nil
	}

	return img, format, errors.New("Storage type unknown")
}

func ReadImageBytes(fname string) ([]byte, error) {
	cfg := config.GetConfig()

	var data []byte

	if cfg.Storage.IsS3() {
		data, err := s3store.ReadImageBytes(fname)

		if err != nil {
			return data, err
		} else {
			return data, nil
		}
	}

	if cfg.Storage.IsLocal() {

		data, err := local.ReadImageBytes(fname)

		if err != nil {
			return data, err
		} else {
			return data, nil
		}
	}

	return data, errors.New("Storage type unknown")
}
