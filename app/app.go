package app

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	goRedis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ShelbyKS/Roamly-backend/app/config"
	_ "github.com/ShelbyKS/Roamly-backend/docs"
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage/postgresql"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage/redis"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi"
	"github.com/ShelbyKS/Roamly-backend/pkg/scheduler"
)

type Roamly struct {
	config  *config.Config
	logger  *logrus.Logger
	pgDB    *gorm.DB
	redisDB *goRedis.Client
}

func New(cfg *config.Config, lg *logrus.Logger) *Roamly {
	return &Roamly{
		config: cfg,
		logger: lg,
	}
}

func (app *Roamly) Run() {
	app.initDBs()

	r := app.newRouter()

	app.initExternalClients()
	app.initAPI(r)

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

func (app *Roamly) initDBs() {
	gormLogger := logger.New(
		app.logger, logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      true,
		},
	)

	pgDB, err := gorm.Open(postgres.Open(app.config.GetPostgresCfg()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	redisClient := goRedis.NewClient(&goRedis.Options{
		Addr:     app.config.Redis.Host + ":" + app.config.Redis.Port,
		Password: app.config.Redis.Password,
	})
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	err = pgDB.AutoMigrate(&orm.User{}, &orm.Trip{}, &orm.Place{})
	if err != nil {
		log.Fatalf("Failed to migrate db: %v", err)
	}

	app.pgDB = pgDB
	app.redisDB = redisClient
}

func (app *Roamly) initAPI(router *gin.Engine) {
	userStorage := postgresql.NewUserStorage(app.pgDB)
	sessionStorage := redis.NewSessionStorage(app.redisDB)
	tripStorage := postgresql.NewTripStorage(app.pgDB)
	placeStorage := postgresql.NewPlaceStorage(app.pgDB)
	mw := middleware.InitMiddleware(sessionStorage)

	schedulerCLient := scheduler.NewClient(scheduler.URL) //todo: move to external

	schedulerService := service.NewShedulerService(schedulerCLient)
	userService := service.NewUserService(userStorage, sessionStorage)
	authService := service.NewAuthService(userStorage, sessionStorage)
	tripService := service.NewTripService(tripStorage, placeStorage)
	placeService := service.NewPlaceService(placeStorage, tripStorage)
	eventService := service.NewEventService(eventStorage, tripStorage, placeStorage)

	handler.NewAuthHandler(router, app.logger, authService, mw)
	router.Use(mw.AuthMiddleware())
	handler.NewUserHandler(router, app.logger, userService)
	handler.NewTripHandler(router, app.logger, tripService, schedulerService)
	handler.NewPlaceHandler(router, app.logger, placeService)
	handler.NewEventHandler(router, app.logger, eventService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (app *Roamly) initExternalClients() {
	googleapi.Init(app.config.GoogleApiKey)
}
