package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/model"
	"github.com/erpe/image_service_go/app/storage"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"image"
	"log"
	"net/http"
	"strconv"
)

var appConfig *config.Config

func init() {
	appConfig = config.GetConfig()
}

/* GET /api/images */
func GetImages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	/*
		optional query argument 'client'
		to scope images by client
	*/
	client := r.FormValue("client")

	images := []model.Image{}

	if len(client) > 0 {
		log.Println("Query scoped by client: ", client)
		if err := db.Preload("Variants").
			Find(&images, "client = ?", client).
			Error; err != nil {
			log.Fatal("ERROR: ", err)
		}
	} else {
		log.Println("No client scope present")
		if err := db.Preload("Variants").
			Find(&images).
			Error; err != nil {
			log.Fatal("ERROR: ", err)
		}
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, images)
}

/* GET /api/images/1 */
func GetImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId := makeInt(vars["id"])

	image := getImageOr404(db, imageId, w)

	if image == nil {
		return
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, image)
}

/* POST /api/images */
func CreateImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	postImage := model.PostImage{}

	json.NewDecoder(r.Body).Decode(&postImage)

	defer r.Body.Close()

	imgBytes, format, err := postImage.Bytes()

	if err != nil {
		log.Println("err: ", err.Error())
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	mErr := supplyMeta(&postImage)

	if mErr != nil {
		respondError(w, http.StatusInternalServerError, mErr.Error())
		return
	}

	if err := db.Save(&postImage).Error; err != nil {

		respondError(w, http.StatusInternalServerError, err.Error())

	} else {

		img := model.Image{}

		db.First(&img, postImage.ID)

		fname := createImageName(postImage.ID, format)

		url, err := storage.SaveImage(imgBytes, fname)

		if err != nil {
			log.Println("ERROR: ", err.Error())
			respondError(w, http.StatusInternalServerError, err.Error())
		}

		img.Url = url
		img.Filename = fname

		if err := db.Save(&img).Error; err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
		} else {
			respondJSON(w, http.StatusCreated, img)
		}
	}
}

func UpdateImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	imageId := makeInt(vars["id"])

	image := getImageOr404(db, imageId, w)

	if image == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	if err := decoder.Decode(&image); err != nil {
		log.Println("ERROR: ", err.Error())
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.Save(&image).Error; err != nil {
		log.Println("ERROR: ", err.Error())
		respondError(w, http.StatusInternalServerError, err.Error())
	}

	respondJSON(w, http.StatusOK, image)
}

func DestroyImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	imageId := makeInt(vars["id"])

	image := getImageOr404(db, imageId, w)

	if image == nil {
		return
	}

	if len(image.Variants) > 0 {
		err := errors.New("Remove existing variants first")
		respondError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err := storage.UnlinkImage(image.Filename); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		db.Unscoped().Delete(&image)
		respondJSON(w, http.StatusNoContent, nil)
	}
}

func createImageName(id int, format string) string {
	return strconv.Itoa(id) + "." + format
}

func getImageOr404(db *gorm.DB, id int, w http.ResponseWriter) *model.Image {
	image := model.Image{}
	if err := db.Preload("Variants").First(&image, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &image
}

func supplyMeta(pImg *model.PostImage) error {

	data, _, err := pImg.Bytes()

	if err != nil {
		return err
	}

	imgReader := bytes.NewReader(data)

	cfg, fmt, err := image.DecodeConfig(imgReader)

	if err != nil {
		return err
	}

	pImg.Width = cfg.Width
	pImg.Height = cfg.Height
	pImg.Format = fmt

	return nil
}
