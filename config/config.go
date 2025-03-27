package config

import (
	"encoding/gob"
	"log"
	"os"
)

type BrewfatherConfig struct {
	AuthToken string
}

type AuthConfig struct {
	Issuer       			string
	ClientID     			string
	ClientSecret 			string
	RedirectURI  			string
	LogoutURI 				string
	PostLogoutRedirectURI 	string
}

var (
	Brewfather BrewfatherConfig
	Authentik  AuthConfig
)

func Init() {
	// Configure encoding
	gob.Register(map[string]interface{}{})

	// Auth token for Brewfather
	authToken := os.Getenv("NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN")
	if authToken == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN is not set")
	}

	// Config object with Brewfather values
	Brewfather = BrewfatherConfig{
		AuthToken: authToken,
	}

	log.Println("Successfully loaded Brewfather config")

	// Authentik auth issuer
	issuer := os.Getenv("NO_GEEKS_BREWING_AUTHENTIK_ISSUER_URL")
	if issuer == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_AUTHENTIK_ISSUER_URL is not set")
	}

	// Authentik client id
	clientID := os.Getenv("NO_GEEKS_BREWING_AUTHENTIK_CLIENT_ID")
	if clientID == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_AUTHENTIK_CLIENT_ID is not set")
	}

	clientSecret := os.Getenv("NO_GEEKS_BREWING_AUTHENTIK_CLIENT_SECRET")
	if clientSecret == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_AUTHENTIK_CLIENT_SECRET is not set")
	}

	// Authentik redirect URL
	redirectURI := os.Getenv("NO_GEEKS_BREWING_REDIRECT_URI")
	if redirectURI == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_REDIRECT_URI is not set")
	}

	// Authentik redirect URL
	logoutURI := os.Getenv("NO_GEEKS_BREWING_AUTHENTIK_LOGOUT_URI")
	if logoutURI == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_AUTHENTIK_LOGOUT_URI is not set")
	}

	// Authentik redirect URL
	postLogoutRedirectURI := os.Getenv("NO_GEEKS_BREWING_POST_LOGOUT_REDIRECT_URI")
	if postLogoutRedirectURI == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_POST_LOGOUT_REDIRECT_URI is not set")
	}

	Authentik = AuthConfig{
		Issuer:       issuer,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		LogoutURI: logoutURI,
		PostLogoutRedirectURI: postLogoutRedirectURI,
	}

	log.Println("Successfully loaded Authentik config")
}
