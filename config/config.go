package config

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	PortalOauthConfig *oauth2.Config
	AdminOauthConfig  *oauth2.Config
)

// GetOauthConfig returns the OAuth config matching the given portal identifier.
// Use "admin" for the admin portal, anything else returns the user portal config.
func GetOauthConfig(portal string) *oauth2.Config {
	if portal == "admin" {
		return AdminOauthConfig
	}
	return PortalOauthConfig
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	PortalOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL_PORTAL"),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint:     google.Endpoint,
	}

	AdminOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL_ADMIN"),
		Scopes:       []string{"email", "profile", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func NewSessionStore(secret string, maxAge int) cookie.Store {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
	})
	return store
}
