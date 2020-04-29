package handler

import (
	"encoding/json"
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/model"
	"github.com/erpe/image_service_go/app/storage"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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

	postImage := model.PostImage{}

	json.NewDecoder(r.Body).Decode(&postImage)

	defer r.Body.Close()

	imgBytes, format, err := postImage.Bytes()

	if err != nil {
		log.Println("err: ", err.Error())
		respondError(w, http.StatusInternalServerError, err.Error())
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

func DestroyImage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	imageId := makeInt(vars["id"])

	image := getImageOr404(db, imageId, w, r)

	if image == nil {
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
	return strconv.Itoa(id) + "-upload." + format
}

func getImageOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Image {
	image := model.Image{}
	if err := db.First(&image, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}

	return &image
}
