package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type Config struct {
	Gmail struct {
		Credentials string `json:"credentials"`
		Token       string `json:"token"`
		From        string `json:"from"`
	} `json:"gmail"`
}

func main() {
	cfgStr := os.Getenv("GMAIL_CREDENTIALS")
	if cfgStr == "" {
		log.Fatalf("Error: GMAIL_CREDENTIALS environment variable is not set")
	}

	var cfg Config
	if err := json.Unmarshal([]byte(cfgStr), &cfg); err != nil {
		log.Fatalf("Error parsing credentials: %v", err)
	}

	// Parse token
	var token oauth2.Token
	if err := json.Unmarshal([]byte(cfg.Gmail.Token), &token); err != nil {
		log.Fatalf("Error parsing token: %v", err)
	}

	// Create OAuth2 config
	oauthConfig, err := google.ConfigFromJSON([]byte(cfg.Gmail.Credentials), gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}

	// Create Gmail service
	client := oauthConfig.Client(context.Background(), &token)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail client: %v", err)
	}

	// Create email message
	to := "masliankovladik@gmail.com"
	subject := "Test Email from Weather API"
	body := "This is a test email from the Weather API service."

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", cfg.Gmail.From, to, subject, body)
	msg := []byte(message)

	// Send email
	_, err = srv.Users.Messages.Send("me", &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(msg),
	}).Do()

	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
