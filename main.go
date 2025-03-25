package main

import (
	"log"
	"net/http"
	"os"
	// "sort"

	"github.com/allan-lewis/no-geeks-brewing-go/batches"
	"github.com/allan-lewis/no-geeks-brewing-go/index"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var batchesMap = make(map[string]batches.Batch)


func main() {
	authToken := os.Getenv("NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN")

	// Check if the environment variable is set
	if authToken == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN is not set")
	}

	ch := make(chan batches.Batch)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Define routes
	r.Get("/", index.IndexHandler) // Serve the index page at "/"

	go func() {
		for {
			// Receive the batch from the channel (this will block until something is received)
			receivedBatch := <-ch
			// Process the received batch (e.g., print it)
			log.Printf("Received batch: %+v", receivedBatch)

			batchesMap[receivedBatch.ID] = receivedBatch
		}
	}()

	go batches.KickOff(ch, authToken)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
