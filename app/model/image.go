package model

import (
	"encoding/base64"
	_ "golang.org/x/image/tiff"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
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
