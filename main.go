package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/allan-lewis/no-geeks-brewing-go/batches"
	"github.com/allan-lewis/no-geeks-brewing-go/index"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func init() {
	gob.Register(map[string]interface{}{})
}

var batchesMap = make(map[string]batches.Batch)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

var (
	oauth2Config *oauth2.Config
	provider     *oidc.Provider
)

func main() {
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Change to true in production with HTTPS
		MaxAge:   3600, // Session expiry in seconds
	}
	
	// Auth token for Brewfather
	authToken := os.Getenv("NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN")
	if authToken == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_BREWFATHER_AUTH_TOKEN is not set")
	}

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
	redirectURL := os.Getenv("NO_GEEKS_BREWING_REDIRECT_URI")
	if redirectURL == "" {
		// Exit the program if the environment variable is not found
		log.Fatal("Environment variable NO_GEEKS_BREWING_REDIRECT_URI is not set")
	}

	var err error
	provider, err = oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     provider.Endpoint(),
	}

	// oauth2Config = &oauth2.Config{
	// 	ClientID:     clientID,
	// 	ClientSecret: clientSecret,
	// 	RedirectURL:  redirectURL,
	// 	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	// 	Endpoint:     provider.Endpoint(),
	// }

	ch := make(chan batches.Batch)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Define routes
	r.Get("/", homeHandler) // Serve the index page at "/"

	// r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, oauth2Config.AuthCodeURL("state123"), http.StatusFound)
	// })

	r.Get("/login", loginHandler)

	r.Get("/auth/callback/authentik", authCallbackHandler)

	go func() {
		for {
			// Receive the batch from the channel (this will block until something is received)
			receivedBatch := <-ch
			// Process the received batch (e.g., print it)
			log.Printf("Received batch: %+v", receivedBatch)

			batchesMap[receivedBatch.ID] = receivedBatch
		}
	}()

	// go batches.KickOff(ch, authToken)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "auth-session")
	userData, ok := session.Values["user"].(map[string]interface{})

	log.Printf("User data %v %v", userData, ok)

	var values []batches.Batch

	err := index.Index(batches.Batches(values)).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := "random-state" // Generate this securely in production
	url := oauth2Config.AuthCodeURL(state)

	// Return a simple HTML response indicating the redirect
	fmt.Fprintf(w, `
	<p>Redirecting to login...</p>
	<script>window.location.href = "%s";</script>
	`, url)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in callback", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch user info
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse user claims
	var claims map[string]interface{}
	if err := userInfo.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save user info in session
	session, _ := store.Get(r, "auth-session")
	session.Values["user"] = claims
	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userData, ok := session.Values["user"].(map[string]interface{})
	log.Printf("User data >> %v %v", userData, ok)

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
}
