package main

import (
	"github.com/ShelbyKS/Roamly-backend/app"
	"log"
)

func main() {
	// load config

	//init logger

	application, err := app.New()
	if err != nil {
		log.Fatal("Failed to create application: %v", err)
	}

	application.Run()
}
