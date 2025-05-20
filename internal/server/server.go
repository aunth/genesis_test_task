package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"weather-service/internal/handlers"
)

func NewServer() *gin.Engine {
	router := gin.New()

	// Logging middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		log.Printf("%s %s | %d | %s", c.Request.Method, c.Request.URL.Path, status, latency)
	})

	// Serve static files
	router.Static("/static", "./internal/static")

	// Serve index.html at root
	router.GET("/", func(c *gin.Context) {
		c.File("./internal/static/index.html")
	})

	// API routes
	router.GET("/weather", handlers.GetWeather)
	router.POST("/subscribe", handlers.Subscribe)
	router.GET("/confirm/:token", handlers.Confirm)
	router.POST("/unsubscribe", handlers.Unsubscribe)

	return router
}

func Start(router *gin.Engine, port string) error {
	return router.Run(":" + port)
}
