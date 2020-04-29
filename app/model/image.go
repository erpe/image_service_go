package model

import (
	"encoding/base64"
	"errors"
	"image"
	"log"
	"net/http"
	"strings"
)

type Image struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Filename  string `json:"filename"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	Variants  []Variant
}

/* get a grip on post-data for handler.CreateImage */
type PostImage struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Filename  string `json:"filename"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	// ignore Data while storing to db
	Data string `gorm:"-" json:"data"`
}

func (pi *PostImage) TableName() string {
	return "images"
}

func (pi *PostImage) Image() (image.Image, string, error) {

	reader := strings.NewReader(pi.Data)

	decoded := base64.NewDecoder(base64.StdEncoding, reader)

	img, format, err := image.Decode(decoded)

	return img, format, err
}

/* returns []byte image data, imagetype, error */
func (pi *PostImage) Bytes() ([]byte, string, error) {

	b, err := base64.StdEncoding.DecodeString(pi.Data)

	if err != nil {
		return []byte(""), "", err
	}

	imgType, err := getImageType(b)

	if err != nil {
		return b, imgType, err
	}

	log.Println("Format: ", imgType)

	return b, imgType, nil
}

func getImageType(b []byte) (string, error) {

	allowed := []string{"jpeg", "jpg", "gif", "png", "tif", "tiff"}

	arr := strings.Split(http.DetectContentType(b), "/")

	if len(arr) < 2 {
		return "", errors.New("Unknown Content-Type:" + arr[0])
	}

	imageType := arr[len(arr)-1]

	for _, item := range allowed {
		if item == imageType {
			return item, nil
		}
	}
	return "", errors.New("Unregistered ImageType: " + imageType)
}
