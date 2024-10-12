package app

import (
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ShelbyKS/Roamly-backend/app/config"
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi"
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
	// pgDB := &gorm.DB{}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s", // sslmode=%s",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"postgres",
		// conf.User,
		// conf.DBName,
		// conf.Password,
		// conf.SSLMode,
	)

	pgDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}

	pgDB.AutoMigrate(&orm.User{}, &orm.Trip{})
	// pgDB.AutoMigrate(&model.Trip{})

	r := app.newRouter()

	app.initExternalClients()
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
	userStorage := storage.NewUserStorage(postgres)
	userService := service.NewUserService(userStorage)
	handler.NewUserHandler(router, app.logger, userService)

	tripStorage := storage.NewTripStorage(postgres)
	tripService := service.NewTripService(tripStorage)
	handler.NewTripHandler(router, app.logger, tripService)
	handler.NewUserHandler(router, app.logger, userService, googleapi.DefaultClient)
}

func (app *Roamly) initExternalClients() {
	googleapi.Init(app.config.GoogleApiKey)
}
