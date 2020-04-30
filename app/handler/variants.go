package handler

import (
	"github.com/erpe/image_service_go/app/model"
	//"github.com/erpe/image_service_go/app/storage"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

/* GET /api/variants */
func GetVariants(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	variants := []model.Variant{}

	varId, ok := vars["imageId"]

	if ok {
		id := makeInt(varId)
		log.Printf("vars:%s", id)

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
	// TODO
}

func DestroyVariant(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// TODO
}

func getVariantOr404(db *gorm.DB, id int, w http.ResponseWriter) *model.Variant {
	variant := model.Variant{}

	if err := db.First(&variant, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &variant
}
