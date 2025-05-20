package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailService struct {
	service *gmail.Service
	from    string
}

func NewEmailService() (*EmailService, error) {
	requiredVars := []string{"GMAIL_CREDENTIALS", "GMAIL_TOKEN", "GMAIL_FROM"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", v)
		}
	}

	credentials := os.Getenv("GMAIL_CREDENTIALS")
	config, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GMAIL_CREDENTIALS: %v", err)
	}

	tokenStr := os.Getenv("GMAIL_TOKEN")
	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		return nil, fmt.Errorf("failed to parse GMAIL_TOKEN: %v", err)
	}
	tokenSource := config.TokenSource(context.Background(), &token)

	client := oauth2.NewClient(context.Background(), tokenSource)
	service, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gmail service: %v", err)
	}

	from := os.Getenv("GMAIL_FROM")
	return &EmailService{
		service: service,
		from:    from,
	}, nil
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	if s == nil || s.service == nil {
		return fmt.Errorf("email service not properly initialized")
	}
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", s.from, to, subject, body)
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(message))
	gmailMessage := &gmail.Message{Raw: encodedMessage}
	_, err := s.service.Users.Messages.Send("me", gmailMessage).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func (s *EmailService) SendConfirmationEmail(to, token string) error {
	if s == nil || s.service == nil {
		return fmt.Errorf("email service not properly initialized")
	}
	confirmationLink := fmt.Sprintf("http://localhost:8080/confirm/%s", token)
	body := fmt.Sprintf("<h2>Welcome to Weather Updates!</h2>"+
		"<p>Please confirm your subscription by clicking the link below:</p>"+
		"<p><a href=\"%s\">Confirm Subscription</a></p>"+
		"<p>If you didn't request this subscription, please ignore this email.</p>", confirmationLink)
	return s.SendEmail(to, "Confirm your weather subscription", body)
}
