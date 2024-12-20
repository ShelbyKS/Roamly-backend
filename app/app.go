package app

import (
	"context"
	"github.com/ShelbyKS/Roamly-backend/internal/utils"
	"log"
	"time"

	"github.com/ShelbyKS/Roamly-backend/internal/middleware"
	"gorm.io/gorm/logger"

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
	"github.com/ShelbyKS/Roamly-backend/pkg/chatgpt"
	"github.com/ShelbyKS/Roamly-backend/pkg/googleapi"
	"github.com/ShelbyKS/Roamly-backend/pkg/kafka"
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

	if err := r.Run(":" + app.config.ServerPort); err != nil {
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

	err = pgDB.AutoMigrate(&orm.User{}, &orm.Trip{}, &orm.Place{}, &orm.Event{},
		&orm.TripUsers{}, &orm.Invite{}, orm.AIChatMessage{})

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
	eventStorage := postgresql.NewEventStorage(app.pgDB)
	inviteStorage := postgresql.NewInviteStorage(app.pgDB)
	aiChatStorage := postgresql.NewAIChatStorage(app.pgDB)

	openAIClient := chatgpt.NewChatGPTClient(app.config.OpenAiKey) //todo: move to external
	// if err != nil {
	// 	log.Fatalf("Failed to create chat-gpt-client %v", err)
	// }
	googleApi := googleapi.NewClient(app.config.GoogleApiKey) //todo: move to external

	producer := kafka.NewMessageBrokerProducer(app.config.Kafka.Host, app.config.Kafka.Port, app.config.Kafka.Topic)
	notifyUrils := utils.NewNotifyUtils(tripStorage, sessionStorage, producer)

	schedulerService := service.NewShedulerService(openAIClient, googleApi, tripStorage, eventStorage, placeStorage, sessionStorage, producer)
	userService := service.NewUserService(userStorage, sessionStorage)
	authService := service.NewAuthService(userStorage, sessionStorage)
	tripService := service.NewTripService(tripStorage, placeStorage, googleApi, openAIClient, sessionStorage, producer, aiChatStorage)
	placeService := service.NewPlaceService(placeStorage, tripStorage, googleApi, eventStorage, openAIClient, sessionStorage, producer)
	eventService := service.NewEventService(eventStorage, tripStorage, placeStorage, sessionStorage, producer)
	inviteService := service.NewInviteService(inviteStorage, tripStorage, app.config.JWTSecret)
	aiChatService := service.NewAIChatService(aiChatStorage, tripStorage, sessionStorage, notifyUrils, openAIClient, googleApi)

	middleware.Mw = middleware.InitMiddleware(sessionStorage)
	router.Use(middleware.Mw.CORSMiddleware())

	handler.NewAuthHandler(router, app.logger, authService)
	handler.NewUserHandler(router, app.logger, userService)
	handler.NewTripHandler(router, app.logger, tripService, placeService, schedulerService)
	handler.NewPlaceHandler(router, app.logger, placeService, *googleApi)
	handler.NewEventHandler(router, app.logger, eventService, tripService)
	handler.NewInviteHandler(router, app.logger, inviteService, tripService)
	handler.NewAIChatHandler(router, app.logger, aiChatService, tripService)

	router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (app *Roamly) initExternalClients() {
	googleapi.Init(app.config.GoogleApiKey)
}
