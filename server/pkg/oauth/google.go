package oauth

import (
	"oauth/internal/config"
	
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
)

var GoogleOAuthConfig *oauth2.Config 

func InitGoogleOAuth() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID: config.Config.GithubClientID,
		ClientSecret: config.Config.GithubClientSecret,
		RedirectURL: config.Config.GithubCallbackURL,

		Scopes: []string{
			"https://www/googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}