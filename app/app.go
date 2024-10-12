package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/ShelbyKS/Roamly-backend/app/config"
	"github.com/ShelbyKS/Roamly-backend/internal/database/orm"
	"github.com/ShelbyKS/Roamly-backend/internal/database/storage"
	"github.com/ShelbyKS/Roamly-backend/internal/handler"
	"github.com/ShelbyKS/Roamly-backend/internal/service"
	"github.com/ShelbyKS/Roamly-backend/pkg/scheduler"
	"gorm.io/driver/postgres"
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

	// pgDB.AutoMigrate(&orm.Trip{})
	err = pgDB.AutoMigrate(&orm.User{}, &orm.Trip{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

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

// type StubScedulerService struct {
// }

// func (*StubScedulerService) GetSchedule(ctx context.Context, places []model.Place) (model.Schedule, error) {
// 	return model.Schedule{}, nil
// }

// type StubScedulerService struct {
// }

func (app *Roamly) initAPI(router *gin.Engine, postgres *gorm.DB) {
	userStorage := storage.NewUserStorage(postgres)
	userService := service.NewUserService(userStorage)
	handler.NewUserHandler(router, app.logger, userService)

	tripStorage := storage.NewTripStorage(postgres)
	tripService := service.NewTripService(tripStorage)

	schedulerCLient := scheduler.NewClient(scheduler.URL)
	schedulerService := service.NewShedulerService(schedulerCLient)
	// log.Println("schedulerService", schedulerService)
	handler.NewTripHandler(router, app.logger, tripService, schedulerService)
}
