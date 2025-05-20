package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func main() {
	credentials := os.Getenv("GMAIL_CREDENTIALS")
	if credentials == "" {
		log.Fatalf("Error: Required environment variable GMAIL_CREDENTIALS is not set")
	}

	oauthConfig, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file: %v", err)
	}

	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
	fmt.Print("Enter the authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	// Exchange auth code for token
	tok, err := oauthConfig.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}

	// Print token
	tokenJSON, err := json.Marshal(tok)
	if err != nil {
		log.Fatalf("Unable to marshal token: %v", err)
	}
	fmt.Printf("Add this to your environment variables as GMAIL_TOKEN:\n%s\n", string(tokenJSON))
}
