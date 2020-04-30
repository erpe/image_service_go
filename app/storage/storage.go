package storage

import (
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/storage/local"
	"github.com/erpe/image_service_go/app/storage/s3store"
)

func SaveImage(data []byte, name string) (string, error) {
	cfg := config.GetConfig()

	var url string

	if cfg.Storage.IsS3() {
		res, err := s3store.SaveImage(data, name)

		if err != nil {
			return "", err
		}
		url = res
	}

	if cfg.Storage.IsLocal() {
		res, err := local.SaveImage(data, name)

		if err != nil {
			return "", err
		}

		url = res
	}

	return url, nil
}

func UnlinkImage(fname string) error {
	cfg := config.GetConfig()
	if cfg.Storage.IsS3() {
		return s3store.UnlinkImage(fname)
	}

	if cfg.Storage.IsLocal() {
		return local.UnlinkImage(fname)
	}
	return nil
}
