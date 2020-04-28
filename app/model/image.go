package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"
)

type Image struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	Variants  []Variant
}

/* get a grip on post-data for handler.CreateImage */
type PostImage struct {
	ID        int    `gorm:"column:id" json:"id"`
	Url       string `json:"url"`
	Alt       string `json:"alt"`
	Copyright string `json:"copyright"`
	Category  string `gorm:"INDEX" json:"category"`
	// ignore Data while storing to db
	Data string `gorm:"-" json:"data"`
}

func (pi *PostImage) TableName() string {
	return "images"
}

func (pi *PostImage) ToImage() (image.Image, string, error) {

	reader := strings.NewReader(pi.Data)

	decoded := base64.NewDecoder(base64.StdEncoding, reader)

	img, format, err := image.Decode(decoded)

	return img, format, err
}

/* returns []byte image data, format, error */
func (pi *PostImage) Bytes() ([]byte, string, error) {

	img, format, err := pi.ToImage()

	if err != nil {
		return []byte(""), "", err
	}

	buf := new(bytes.Buffer)

	var imgErr error

	switch format {
	case "jpeg":
		imgErr = jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
	case "png":
		imgErr = png.Encode(buf, img)
	case "gif":
		imgErr = gif.Encode(buf, img, nil)
	case "tif":
		imgErr = tiff.Encode(buf, img, nil)
	default:
		imgErr = errors.New("Unsupported format: " + format)
	}

	if imgErr != nil {
		return []byte(""), format, imgErr
	}

	return buf.Bytes(), format, nil
}
