package router

import (
	"log"
	"net/http"

	"github.com/allan-lewis/no-geeks-brewing-go/layout"
	"github.com/allan-lewis/no-geeks-brewing-go/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// var batchesMap = make(map[string]batches.Batch)

func Start() {
	// ch := make(chan batches.Batch)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Define routes
	// Serve the index page at "/"

	// r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, oauth2Config.AuthCodeURL("state123"), http.StatusFound)
	// })

	// Main index/layout route
	r.Get("/", layout.LayoutHandler)

	// OAuth routes
	r.Get("/auth/login", oauth.LoginHandler)
	r.Get("/auth/callback/authentik", oauth.AuthCallbackHandler)
	r.Get("/auth/logout", oauth.LogoutHandler)

	// go func() {
	// 	for {
	// 		// Receive the batch from the channel (this will block until something is received)
	// 		receivedBatch := <-ch
	// 		// Process the received batch (e.g., print it)
	// 		log.Printf("Received batch: %+v", receivedBatch)

	// 		batchesMap[receivedBatch.ID] = receivedBatch
	// 	}
	// }()

	// go batches.KickOff(ch, authToken)

	// Start the HTTP server
	log.Println("HTTP server is listening on port 8080")
	http.ListenAndServe(":8080", r)
}
