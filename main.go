package main

import (
	"log"
	"net/http"
	"os"

	"github.com/allan-lewis/no-geeks-brewing-go/batch"
	"github.com/allan-lewis/no-geeks-brewing-go/templates"
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var batchesMap = make(map[string]batch.Batch)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins during development
		return true
	},
}

// Handle WebSocket connections
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		if _, _, err := conn.NextReader(); err != nil {
			log.Println("WebSocket connection closed:", err)
			delete(clients, conn)
			break
		}
	}
}

// Watch for file changes
func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating file watcher:", err)
	}
	defer watcher.Close()

	// Add directories to watch
	directories := []string{"./"}
	for _, dir := range directories {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatalf("Error watching directory %s: %v", dir, err)
		}
	}

	log.Println("Watching for changes...")
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				log.Printf("Change detected: %s", event.Name)
				notifyClients()
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}

// Notify connected clients to reload
func notifyClients() {
	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			log.Println("Error sending reload message:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

// Serve the index page with embedded JS for live reload
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Create a slice to hold the values from the map
	var values []batch.Batch

	// Iterate over the map and append the values to the slice
	for _, value := range batchesMap {
		values = append(values, value)
	}

	err := templates.Index(templates.Batches(values)).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func main() {
	authToken := os.Getenv("NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN")

	// Check if the environment variable is set
	if authToken == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN is not set")
	}

	ch := make(chan batch.Batch)

	r := chi.NewRouter()

	// Define routes
	r.Get("/", indexHandler) // Serve the index page at "/"
	r.Get("/ws", wsHandler)  // WebSocket connection at "/ws"

	http.HandleFunc("/ws", wsHandler)
	go func() {
		log.Println("WebSocket server started on ws://localhost:9090/ws")
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()

	go func() {
		for {
			// Receive the batch from the channel (this will block until something is received)
			receivedBatch := <-ch
			// Process the received batch (e.g., print it)
			log.Printf("Received batch: %+v", receivedBatch)

			batchesMap[receivedBatch.ID] = receivedBatch
		}
	}()

	// Start file watcher
	go watchFiles()

	go batch.KickOff(ch, authToken)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
