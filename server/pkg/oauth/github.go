package oauth

import (
	"oauth/internal/config"

	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

var GithubOAuthConfig *oauth2.Config

func InitGithubOAuth() {
	GithubOAuthConfig = &oauth2.Config{
		ClientID:     config.Config.GithubClientID,
		ClientSecret: config.Config.GithubClientSecret,
		RedirectURL:  config.Config.GithubCallbackURL,
		Scopes: []string{
			"read:user",
			"user:email",
		},
		Endpoint: githubOAuth.Endpoint,
	}
}
