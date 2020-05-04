package app

import (
	"fmt"
	"github.com/erpe/image_service_go/app/config"
	"github.com/erpe/image_service_go/app/handler"
	"github.com/erpe/image_service_go/app/model"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Config *config.Config
}

func (a *App) Initialize(config *config.Config) {

	a.Config = config

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

	a.setRouters()
	a.setMiddleware()
}

func (a *App) setMiddleware() {
	a.Router.Use(handler.LoggingMiddleware)
}

func (a *App) setRouters() {
	a.Router = mux.NewRouter()
	apiRouter := a.Router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/images/{imageId}/variants", a.variantsHandler).
		Methods("GET")
	apiRouter.HandleFunc("/images/{imageId}/variants", a.createVariantHandler).
		Methods("POST")

	apiRouter.HandleFunc("/images", a.imagesHandler).
		Methods("GET")
	apiRouter.HandleFunc("/images/{id}", a.imageHandler).
		Methods("GET")
	apiRouter.HandleFunc("/images", a.createImageHandler).
		Methods("POST")
	apiRouter.HandleFunc("/images/{id}", a.destroyImageHandler).
		Methods("DELETE")
	apiRouter.HandleFunc("/images/{id}", a.updateImageHandler).
		Methods("PATCH")

	apiRouter.HandleFunc("/variants", a.variantsHandler).
		Methods("GET")
	apiRouter.HandleFunc("/variants/{id}", a.variantHandler).
		Methods("GET")
	apiRouter.HandleFunc("/variants/{id}", a.destroyVariantHandler).
		Methods("DELETE")
}

func (a *App) imagesHandler(w http.ResponseWriter, r *http.Request) {
	handler.GetImages(a.DB, w, r)
}

func (a *App) imageHandler(w http.ResponseWriter, r *http.Request) {
	handler.GetImage(a.DB, w, r)
}

func (a *App) createImageHandler(w http.ResponseWriter, r *http.Request) {
	handler.CreateImage(a.DB, w, r)
}

func (a *App) destroyImageHandler(w http.ResponseWriter, r *http.Request) {
	handler.DestroyImage(a.DB, w, r)
}

func (a *App) updateImageHandler(w http.ResponseWriter, r *http.Request) {
	handler.UpdateImage(a.DB, w, r)
}

func (a *App) variantsHandler(w http.ResponseWriter, r *http.Request) {
	handler.GetVariants(a.DB, w, r)
}

func (a *App) variantHandler(w http.ResponseWriter, r *http.Request) {
	handler.GetVariant(a.DB, w, r)
}

func (a *App) createVariantHandler(w http.ResponseWriter, r *http.Request) {
	handler.CreateVariant(a.DB, w, r)
}

func (a *App) destroyVariantHandler(w http.ResponseWriter, r *http.Request) {
	handler.DestroyVariant(a.DB, w, r)
}

func (a *App) Run() {
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Authorization",
		"Content-Type",
	})

	originsOk := handlers.AllowedOrigins(a.Config.Server.Cors)

	for _, host := range a.Config.Server.Cors {
		log.Printf("enabling CORS for: %s", host)
	}

	methodsOk := handlers.AllowedMethods([]string{
		"GET",
		"OPTIONS",
		"PUT",
		"POST",
		"PATCH",
		"DELETE",
	})

	log.Printf("listening on '%s'", a.Config.Server.ToString())
	log.Fatal(http.ListenAndServe(a.Config.Server.ToString(),
		handlers.CORS(originsOk, headersOk, methodsOk)(a.Router)))
}
