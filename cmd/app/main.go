package main

import (
	"github.com/ShelbyKS/Roamly-backend/app"
	"github.com/ShelbyKS/Roamly-backend/app/config"
	"github.com/ShelbyKS/Roamly-backend/app/logger"
)

func main() {
	appCfg := config.LoadConfig()

	lg := logger.InitLogger(appCfg)

	application := app.New(appCfg, lg)
	application.Run()
}
