package app

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/app/config"
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi"
	"github.com/ShelbyKS/Roamly-backend/pkg/scheduler"
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
	pgDB, err := gorm.Open(postgres.Open(app.config.GetPostgresCfg()), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = pgDB.AutoMigrate(&orm.User{}, &orm.Trip{}, &orm.Place{})
	if err != nil {
		log.Fatalf("Failed to migrate db: %v", err)
	}

	r := app.newRouter()

	app.initExternalClients()
	app.initAPI(r, pgDB)

	if err := r.Run(":" + string(app.config.ServerPort)); err != nil {
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
	userStorage := storage.NewUserStorage(postgres)
	tripStorage := storage.NewTripStorage(postgres)
	placeStorage := storage.NewPlaceStorage(postgres)

	schedulerCLient := scheduler.NewClient(scheduler.URL) //todo: move to external

	schedulerService := service.NewShedulerService(schedulerCLient)
	userService := service.NewUserService(userStorage)
	tripService := service.NewTripService(tripStorage, placeStorage)
	placeService := service.NewPlaceService(placeStorage, tripStorage)

	handler.NewUserHandler(router, app.logger, userService)
	handler.NewTripHandler(router, app.logger, tripService, schedulerService)
	handler.NewPlaceHandler(router, app.logger, placeService)
}

func (app *Roamly) initExternalClients() {
	googleapi.Init(app.config.GoogleApiKey)
}
