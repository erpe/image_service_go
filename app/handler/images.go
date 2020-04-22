package handler

import (
	"github.com/erpe/image_service_go/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func GetImages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	images := []model.Image{}

	if err := db.Find(&images).Error; err != nil {
		log.Fatal("ERROR: ", err)
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, images)
}

func GetImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId := makeInt(vars["id"])

	image := getImageOr404(db, imageId, w, r)

	if image == nil {
		return
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, image)
}

func getImageOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Image {
	image := model.Image{}
	if err := db.First(&image, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &image
}
