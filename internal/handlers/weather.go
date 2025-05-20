package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"weather-service/internal/models"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	apiKey string
}

func NewWeatherHandler() (*WeatherHandler, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("required environment variable API_KEY is not set")
	}

	return &WeatherHandler{
		apiKey: apiKey,
	}, nil
}

func GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "City is required",
		})
		return
	}

	handler, err := NewWeatherHandler()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to create weather handler: %v", err),
		})
		return
	}

	weather, err := handler.fetchWeatherData(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch weather data: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, weather)
}

func (h *WeatherHandler) fetchWeatherData(city string) (*models.Weather, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, h.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var apiResponse struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity float64 `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		ID int `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	weather := &models.Weather{
		ID:          fmt.Sprintf("%d", apiResponse.ID),
		Temperature: apiResponse.Main.Temp,
		Humidity:    apiResponse.Main.Humidity,
		Description: apiResponse.Weather[0].Description,
	}

	if len(apiResponse.Weather) > 0 {
		weather.Description = apiResponse.Weather[0].Description
	} else {
		weather.Description = "No description available"
	}

	return weather, nil
}
