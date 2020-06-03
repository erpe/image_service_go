package handler

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

type ClientResult struct {
	Client string `json:"client"`
	Count  int    `json:"count"`
}

func GetClients(db *gorm.DB, w http.ResponseWriter, r *http.Request) {

	var clients []ClientResult

	db.Raw("SELECT client, count(images) from images GROUP BY client").Scan(&clients)
	respondJSON(w, http.StatusOK, clients)
}
