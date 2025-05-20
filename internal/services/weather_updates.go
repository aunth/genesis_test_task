package services

import (
	"fmt"
	"log"
	"time"

	"weather-service/internal/database"
)

type WeatherUpdateService struct {
	emailService   *EmailService
	weatherHandler *WeatherHandler
}

func NewWeatherUpdateService(emailService *EmailService, weatherHandler *WeatherHandler) *WeatherUpdateService {
	return &WeatherUpdateService{
		emailService:   emailService,
		weatherHandler: weatherHandler,
	}
}

func (s *WeatherUpdateService) StartScheduler() {
	go func() {
		for {
			// Get all hourly subscriptions
			db, err := database.Connect()
			if err != nil {
				log.Printf("Failed to connect to database: %v", err)
				time.Sleep(5 * time.Minute)
				continue
			}

			rows, err := db.Query(`
				SELECT email, city 
				FROM subscription 
				WHERE confirmed = true AND frequency = 'hourly'`)
			if err != nil {
				log.Printf("Failed to query hourly subscriptions: %v", err)
				db.Close()
				time.Sleep(5 * time.Minute)
				continue
			}

			for rows.Next() {
				var email, city string
				if err := rows.Scan(&email, &city); err != nil {
					log.Printf("Error scanning subscription: %v", err)
					continue
				}

				weather, err := s.weatherHandler.fetchWeatherData(city)
				if err != nil {
					log.Printf("Error fetching weather for %s: %v", city, err)
					continue
				}

				subject := fmt.Sprintf("Hourly Weather Update for %s", city)
				body := fmt.Sprintf(`
					<h2>Hourly Weather Update for %s</h2>
					<p>Current weather conditions:</p>
					<ul>
						<li>Temperature: %.1f°C</li>
						<li>Humidity: %.1f%%</li>
						<li>Conditions: %s</li>
					</ul>
					<p>This update was sent at: %s</p>
					<p>To unsubscribe, click <a href="http://localhost:8080/unsubscribe?token=YOUR_TOKEN">here</a></p>
				`, city, weather.Temperature, weather.Humidity, weather.Description, time.Now().Format(time.RFC1123))

				if err := s.emailService.SendEmail(email, subject, body); err != nil {
					log.Printf("Error sending hourly update to %s: %v", email, err)
					continue
				}

				log.Printf("Sent hourly weather update to %s for %s", email, city)
			}

			rows.Close()
			db.Close()
			time.Sleep(1 * time.Hour)
		}
	}()

	go func() {
		for {
			db, err := database.Connect()
			if err != nil {
				log.Printf("Failed to connect to database: %v", err)
				time.Sleep(5 * time.Minute)
				continue
			}
			defer db.Close()

			rows, err := db.Query(`
				SELECT email, city 
				FROM subscription 
				WHERE confirmed = true AND frequency = 'daily'`)
			if err != nil {
				log.Printf("Failed to query daily subscriptions: %v", err)
				time.Sleep(5 * time.Minute)
				continue
			}

			for rows.Next() {
				var email, city string
				if err := rows.Scan(&email, &city); err != nil {
					log.Printf("Error scanning subscription: %v", err)
					continue
				}

				weather, err := s.weatherHandler.fetchWeatherData(city)
				if err != nil {
					log.Printf("Error fetching weather for %s: %v", city, err)
					continue
				}

				subject := fmt.Sprintf("Daily Weather Update for %s", city)
				body := fmt.Sprintf(`
					<h2>Daily Weather Update for %s</h2>
					<p>Current weather conditions:</p>
					<ul>
						<li>Temperature: %.1f°C</li>
						<li>Humidity: %.1f%%</li>
						<li>Conditions: %s</li>
					</ul>
					<p>This update was sent at: %s</p>
					<p>To unsubscribe, click <a href="http://localhost:8080/unsubscribe?token=YOUR_TOKEN">here</a></p>
				`, city, weather.Temperature, weather.Humidity, weather.Description, time.Now().Format(time.RFC1123))

				if err := s.emailService.SendEmail(email, subject, body); err != nil {
					log.Printf("Error sending daily update to %s: %v", email, err)
					continue
				}

				log.Printf("Sent daily weather update to %s for %s", email, city)
			}

			rows.Close()
			time.Sleep(24 * time.Hour)
		}
	}()
}
