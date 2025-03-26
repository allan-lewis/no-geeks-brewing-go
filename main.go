package main

import (
	"github.com/allan-lewis/no-geeks-brewing-go/config"
	"github.com/allan-lewis/no-geeks-brewing-go/oauth"
	"github.com/allan-lewis/no-geeks-brewing-go/router"
)

func init() {
	// Configure application-wide settings
	config.Init()
}

func main() {
	// Get OAuth ready for use
	oauth.Init(config.Authentik)

	// Start the router (http.ListenAndServe will keep things running)
	router.Start()
}
