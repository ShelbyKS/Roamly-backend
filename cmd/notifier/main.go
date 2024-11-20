package main

import (
	"github.com/ShelbyKS/Roamly-backend/notifier"
	"github.com/ShelbyKS/Roamly-backend/notifier/config"
)

func main() {
	appCfg := config.LoadConfig()

	application := notifier.New(appCfg)
	application.Run()
}
