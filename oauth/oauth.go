package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/allan-lewis/no-geeks-brewing-go/config"
	"github.com/allan-lewis/no-geeks-brewing-go/user"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	oauth2Config          *oauth2.Config
	provider              *oidc.Provider
	logoutURI             string
	postLogoutRedirectURI string
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func Init(authentikConfig config.AuthConfig) {
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Change to true in production with HTTPS
		MaxAge:   3600,  // Session expiry in seconds
	}

	var err error
	provider, err = oidc.NewProvider(context.Background(), config.Authentik.Issuer)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	oauth2Config = &oauth2.Config{
		ClientID:     config.Authentik.ClientID,
		ClientSecret: config.Authentik.ClientSecret,
		RedirectURL:  config.Authentik.RedirectURI,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     provider.Endpoint(),
	}

	logoutURI = authentikConfig.LogoutURI
	postLogoutRedirectURI = authentikConfig.PostLogoutRedirectURI
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	state := "random-state" // Generate this securely in production
	url := oauth2Config.AuthCodeURL(state)

	responseRedirect(w, url)
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
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

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging out")
	// Clear the local session
	session, _ := store.Get(r, "auth-session")
	session.Values = nil
	session.Options.MaxAge = -1 // Delete session
	session.Save(r, w)

	// TODO - Figure out a way to logout of Authentik, as well.
	// This should work per OIDC conventions/standards, but doesn't seem supported
	// out-of-the box for logout flows.
	redirectUrl := fmt.Sprintf("%s?post_logout_redirect_uri=%s&client_id=%s", logoutURI, url.QueryEscape(postLogoutRedirectURI), oauth2Config.ClientID)
	log.Printf("Redirect url %s", redirectUrl)

	responseRedirect(w, "/")
}

func responseRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusNoContent)
}

func User(r *http.Request) user.User {
	session, err := store.Get(r, "auth-session")
	if err != nil {
		return user.UnauthenticatedUser()
	}

	userData, ok := session.Values["user"].(map[string]any)
	if !ok {
		return user.UnauthenticatedUser()
	}

	name, nameOk := userData["name"].(string)
	if !nameOk {
		return user.UnauthenticatedUser()
	}

	email, emailOk := userData["email"].(string)
	if !emailOk {
		return user.UnauthenticatedUser()
	}

	return user.AuthenticatedUser(name, email)
}
