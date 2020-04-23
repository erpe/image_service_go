package handler

import (
	"encoding/json"
	"github.com/erpe/image_service_go/app/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

/* GET /api/images */
func GetImages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	images := []model.Image{}

	if err := db.Find(&images).Error; err != nil {
		log.Fatal("ERROR: ", err)
	}

	defer r.Body.Close()

	respondJSON(w, http.StatusOK, images)
}

/* GET /api/images/1 */
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

/* POST /api/images */
func CreateImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	image := model.PostImage{}

	json.NewDecoder(r.Body).Decode(&image)

	defer r.Body.Close()

	if err := db.Save(&image).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
	} else {
		img := model.Image{}
		db.First(&img, image.ID)
		respondJSON(w, http.StatusCreated, img)
	}
}

func getImageOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Image {
	image := model.Image{}
	if err := db.First(&image, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &image
}
