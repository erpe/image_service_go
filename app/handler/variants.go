package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/erpe/image_service_go/app/model"
	"github.com/erpe/image_service_go/app/storage"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

/* GET /api/variants */
func GetVariants(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	variants := []model.ReadVariant{}

	imgId, ok := vars["imageId"]

	if ok {
		id := makeInt(imgId)

		log.Printf("imgId: %s", id)

		img := getImageOr404(db, id, w)

		if img == nil {
			return
		}

		db.Model(&img).Related(&variants)

	} else {

		log.Println("variants unscoped")

		if err := db.Find(&variants).Error; err != nil {
			log.Fatal("ERROR: ", err.Error())
			respondError(w, http.StatusInternalServerError, err.Error())
		}
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, variants)

}

func GetVariant(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	variantId := makeInt(vars["id"])

	variant := getVariantOr404(db, variantId, w)

	if variant == nil {
		return
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, variant)
}

func CreateVariant(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	imageId := makeInt(vars["imageId"])

	postVar := model.PostVariant{}

	img := getImageOr404(db, imageId, w)

	if img == nil {
		return
	}

	json.NewDecoder(r.Body).Decode(&postVar)

	defer r.Body.Close()

	origin, _, err := img.Image()

	if err != nil {
		log.Println("ERROR - CreateVariant: ", err.Error())
	}

	var newImg image.Image

	width := postVar.Width
	height := postVar.Height

	if postVar.KeepRatio == true {
		newImg = imaging.Resize(origin, width, height, imaging.Lanczos)
	} else {
		newImg = imaging.CropAnchor(origin, width, height, imaging.Center)
	}

	imgBytes, err := encodeImageBytes(newImg, postVar.Filetype)

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	if err := db.Save(&postVar).Error; err != nil {
		log.Println("Save error: ", err.Error())
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	variant := model.Variant{}

	db.First(&variant, postVar.ID)

	fname := createVariantName(imageId, postVar.ID, postVar.Filetype)

	url, err := storage.SaveImage(imgBytes, fname)

	if err != nil {
		log.Println("ERROR storage: ", err.Error())
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	variant.Url = url
	variant.Filename = fname
	variant.ImageID = img.ID
	//variant.Image = *img

	if err := db.Save(&variant).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondJSON(w, http.StatusCreated, variant)
	}
}

func DestroyVariant(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := makeInt(vars["id"])

	variant := getVariantOr404(db, id, w)

	if variant == nil {
		return
	}

	if err := storage.UnlinkImage(variant.Filename); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		db.Unscoped().Delete(&variant)
		respondJSON(w, http.StatusNoContent, nil)
	}

}

func getVariantOr404(db *gorm.DB, id int, w http.ResponseWriter) *model.Variant {
	variant := model.Variant{}

	if err := db.Preload("Image").First(&variant, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &variant
}

func createVariantName(imgId int, varId int, format string) string {
	return strconv.Itoa(imgId) + "-" + strconv.Itoa(varId) + "." + format
}

func encodeImageBytes(img image.Image, format string) ([]byte, error) {

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
