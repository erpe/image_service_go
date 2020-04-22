package app

import (
	"fmt"
	"github.com/austauschkompass/image_service_go/app/config"
	"github.com/austauschkompass/image_service_go/app/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s password=%s",
		config.DB.Host,
		config.DB.Username,
		config.DB.Name,
		config.DB.Password,
	)

	db, err := gorm.Open(config.DB.Dialect, dbURI)

	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}

	a.DB = model.DBMigrate(db)

	if config.Server.Debug == true {
		a.DB.LogMode(true)
	}

	log.Printf("connected database '%s'", config.DB.Name)
	log.Printf("listening on '%s'", config.Server.ToString())
}

func (a *App) Run() {
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Authorization",
		"Content-Type",
	})

	originsOk := handlers.AllowedOrigins([]string{
		"http://localhost:3000",
	})

	methodsOk := handlers.AllowedMethods([]string{
		"GET",
		"OPTIONS",
		"PUT",
		"POST",
		"PATCH",
		"DELETE",
	})

	log.Fatal(http.ListenAndServe("localhost:3000", handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}
