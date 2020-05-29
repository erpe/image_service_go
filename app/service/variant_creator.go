package service

import (
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/erpe/image_service_go/app/model"
	"github.com/erpe/image_service_go/app/storage"
	"github.com/jinzhu/gorm"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
)

type VariantCreator struct {
	DB    *gorm.DB
	Image *model.Image
	// modes -
	// fit: scale down to bounding box
	// fill: resize and crop to fill width/height
	// resize: scale either by width or hight
	Mode string
}

func (vc *VariantCreator) Run(width int, height int, format string, name string) (model.Variant, error) {

	variant := model.Variant{}

	origin, _, err := vc.Image.Image()

	if err != nil {
		return variant, err
	}

	var img image.Image

	//img = imaging.Fill(origin, width, height, imaging.Center, imaging.Lanczos)
	switch vc.Mode {
	case "fit":
		img = imaging.Fit(origin, width, height, imaging.Lanczos)
	case "fill":
		img = imaging.Fill(origin, width, height, imaging.Center, imaging.Lanczos)
	case "resize":
		if width > 0 && height > 0 {
			return variant, errors.New("width and height greater 0 - resize would distort")
		}
		img = imaging.Resize(origin, width, height, imaging.Lanczos)
	default:
		img = imaging.Fit(origin, width, height, imaging.Lanczos)
	}

	variantBytes, err := EncodeImageBytes(img, format)

	if err != nil {
		return variant, err
	}

	// get width, height from variantBytes
	cfg, fmt, err := ExtractMeta(variantBytes)

	if err != nil {
		return variant, err
	}

	// apply width, height, format, name to struct
	variant.Width = cfg.Width
	variant.Height = cfg.Height
	variant.Name = name
	variant.Format = fmt
	variant.ImageID = vc.Image.ID

	if err := vc.DB.Save(&variant).Error; err != nil {
		return variant, err
	}

	// create variant filename
	fname := CreateVariantName(vc.Image.ID, variant.ID, fmt)

	// saveImage bytes, return url
	url, err := storage.SaveImage(variantBytes, fname)

	if err != nil {
		return variant, err
	}

	// save resulting url, filenema to struct
	variant.Url = url
	variant.Filename = fname

	if err := vc.DB.Save(&variant).Error; err != nil {
		return variant, err
	}

	return variant, nil
}

func ExtractMeta(data []byte) (image.Config, string, error) {

	var cfg image.Config

	imgReader := bytes.NewReader(data)
	cfg, fmt, err := image.DecodeConfig(imgReader)

	if err != nil {
		return cfg, "", err
	}

	return cfg, fmt, nil
}

func EncodeImageBytes(img image.Image, format string) ([]byte, error) {

	buf := new(bytes.Buffer)

	var err error

	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
	case "png":
		err = png.Encode(buf, img)
	case "tiff":
		err = tiff.Encode(buf, img, nil)
	case "gif":
		err = gif.Encode(buf, img, nil)
	default:
		err = errors.New("unsupported format: " + format)
	}

	if err != nil {
		return buf.Bytes(), err
	}

	return buf.Bytes(), nil
}

func CreateVariantName(imgId int, varId int, format string) string {
	return strconv.Itoa(imgId) + "-" + strconv.Itoa(varId) + "." + format
}
