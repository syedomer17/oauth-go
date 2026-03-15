package oauth

import (
	"oauth/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     config.Config.GoogleClientID,
		ClientSecret: config.Config.GoogleClientSecret,
		RedirectURL:  config.Config.GoogleCallbackURL,

		Scopes: []string{
			"https://www/googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
