package app

import (
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type Roamly struct {
	//config
	//logger
}

func New() (*Roamly, error) {
	return &Roamly{}, nil
}

func (app *Roamly) Run() {
	//load config

	//init dbs
	//postgres, err := gorm.Open(postgres.Open(app.config.GetDsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("Failed to connect to postgres: %v", err)
	//}
	pgDB := &gorm.DB{}

	r := app.newRouter()

	app.initAPI(r, pgDB)

	//get port from config
	port := "8080"

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (app *Roamly) newRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return router
}

func (app *Roamly) initAPI(router *gin.Engine, postgres *gorm.DB) {
	userStorage := storage.NewStorage(postgres)
	userService := service.NewService(userStorage)
	handler.NewUserHandler(router, userService)
}
