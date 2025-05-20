package main

import (
	"log"
	"os"

	"weather-service/internal/database"
	"weather-service/internal/server"
	"weather-service/internal/services"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := server.NewServer()

	emailService, err := services.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	weatherHandler, err := services.NewWeatherHandler()
	if err != nil {
		log.Fatalf("Failed to initialize weather handler: %v", err)
	}

	weatherUpdateService := services.NewWeatherUpdateService(emailService, weatherHandler)
	weatherUpdateService.StartScheduler()

	_, err = database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Starting weather service on port %s", port)
	if err := server.Start(router, port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
