package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/app/config"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
)

type Roamly struct {
	config *config.Config
	logger *logrus.Logger
}

func New(cfg *config.Config, lg *logrus.Logger) *Roamly {
	return &Roamly{
		config: cfg,
		logger: lg,
	}
}

func (app *Roamly) Run() {
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
	handler.NewUserHandler(router, app.logger, userService)
}
