package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"weather-service/internal/database"
	"weather-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	Email     string `json:"email"`
	City      string `json:"city"`
	Frequency string `json:"frequency"`
}

var emailService *services.EmailService

func init() {
	var err error
	emailService, err = services.NewEmailService()
	if err != nil {
		log.Printf("Failed to initialize email service: %v", err)
	}
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateRequestData(request *SubscriptionRequest) error {
	if request.Email == "" {
		return fmt.Errorf("email is required")
	}

	if !isValidEmail(request.Email) {
		return fmt.Errorf("invalid email")
	}

	if request.City == "" {
		return fmt.Errorf("city is required")
	}

	if request.Frequency == "" {
		return fmt.Errorf("frequency is required")
	}

	if request.Frequency != "hourly" && request.Frequency != "daily" && request.Frequency != "seconds" {
		return fmt.Errorf("invalid frequency")
	}

	return nil
}

func Subscribe(c *gin.Context) {
	// Check required environment variables
	requiredVars := []string{"GMAIL_CREDENTIALS", "GMAIL_TOKEN", "GMAIL_FROM"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Required environment variable %s is not set", v)})
			return
		}
	}

	if emailService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email service not properly initialized"})
		return
	}

	var request SubscriptionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := validateRequestData(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM subscription WHERE email = $1)", request.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already subscribed"})
		return
	}

	token := uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO subscription (email, city, frequency, confirmation_token)
		VALUES ($1, $2, $3, $4)`,
		request.Email, request.City, request.Frequency, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create subscription: %v", err)})
		return
	}

	if err := emailService.SendConfirmationEmail(request.Email, token); err != nil {
		log.Printf("Failed to send confirmation email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})
}

func Confirm(c *gin.Context) {
	fmt.Println("Confirming subscription")
	requiredVars := []string{"GMAIL_CREDENTIALS", "GMAIL_TOKEN", "GMAIL_FROM"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Required environment variable %s is not set", v)})
			return
		}
	}

	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	defer db.Close()

	result, err := db.Exec(`
		UPDATE subscription 
		SET confirmed = true 
		WHERE confirmation_token = $1 AND confirmed = false`,
		token)
	if err != nil {
		log.Println("Error confirming subscription:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}

func Unsubscribe(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing unsubscribe token"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	defer db.Close()

	result, err := db.Exec(`
		DELETE FROM subscription 
		WHERE confirmation_token = $1`,
		token)
	if err != nil {
		log.Println("Error deleting subscription:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscription successful"})
}
