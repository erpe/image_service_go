package storage

import (
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/storage/local"
	"github.com/erpe/image_service_go/app/storage/s3store"
)

func Save(data []byte, name string) (string, error) {
	cfg := config.GetConfig()

	var url string

	if cfg.Storage.IsS3() {
		res, err := s3store.StoreS3(data, name)

		if err != nil {
			return "", err
		}
		url = res
	}

	if cfg.Storage.IsLocal() {
		res, err := local.StoreLocal(data, name)

		if err != nil {
			return "", err
		}

		url = res
	}

	return url, nil
}
