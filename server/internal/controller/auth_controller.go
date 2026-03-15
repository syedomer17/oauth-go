package controller

import (
	"net/http"
	"oauth/internal/models"
	"oauth/internal/services"
	"oauth/internal/utils"
	"oauth/pkg/oauth"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func GoogleLogin(c *gin.Context) {
	state := utils.GenerateState()

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"oauth_state",
		state,
		300,
		"/",
		"",
		true,
		true,
	)

	url := oauth.GoogleOAuthConfig.AuthCodeURL(state)

	c.Redirect(http.StatusFound, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")

	cookieState, err := c.Cookie("oauth_state")

	if err != nil || state == "" || state != cookieState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid oauth state"})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", true, true)

	code := c.Query("code")

	token, err := oauth.GoogleOAuthConfig.Exchange(c, code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token exchange failed",
		})
		return
	}

	client := oauth.GoogleOAuthConfig.Client(c, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch profile",
		})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	json.NewDecoder(resp.Body).Decode(&googleUser)

	authService := services.AuthService{}

	user, err := authService.HandleOAuthLogin(
		googleUser.Email,
		googleUser.Name,
		googleUser.Picture,
		models.ProviderGoogle,
		googleUser.ID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "login failed",
		})
		return
	}

	jwtToken, err := utils.GenerateAccessToken(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	session, err := authService.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create session",
		})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"access_token",
		jwtToken,
		int(utils.AccessTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)
	c.SetCookie(
		"refresh_token",
		session.RefreshToken,
		int(utils.RefreshTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})

}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token is required"})
		return
	}

	authService := services.AuthService{}
	newAccessToken, newRefreshToken, err := authService.RefreshSession(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"access_token",
		newAccessToken,
		int(utils.AccessTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		int(utils.RefreshTokenTTL.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed"})
}

func GithubLogin(c *gin.Context) {
	state := utils.GenerateState()

	// Persist oauth state in a short-lived cookie to validate callback origin.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"oauth_state",
		state,
		300,
		"/",
		"",
		true,
		true,
	)

	url := oauth.GithubOAuthConfig.AuthCodeURL(state)

	c.Redirect(http.StatusFound, url)
}

func GithubCallback(c *gin.Context) {
	state := c.Query("state")

	cookieState, err := c.Cookie("oauth_state")

	// Reject callback when state mismatches to prevent OAuth CSRF attacks.
	if err != nil || state == "" || state != cookieState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid oauth state"})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", true, true)

	code := c.Query("code")

	token, err := oauth.GithubOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token exchange failed",
		})
		return
	}

	client := oauth.GithubOAuthConfig.Client(c, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch github profile",
		})
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID     int64  `json:"id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		Login  string `json:"login"`
		Avatar string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode github profile"})
		return
	}

	if githubUser.Email == "" {
		emailResp, emailErr := client.Get("https://api.github.com/user/emails")
		if emailErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch github email"})
			return
		}
		defer emailResp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode github email"})
			return
		}

		for _, e := range emails {
			if e.Primary && e.Verified {
				githubUser.Email = e.Email
				break
			}
		}

		if githubUser.Email == "" && len(emails) > 0 {
			githubUser.Email = emails[0].Email
		}
	}

	if githubUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "github account has no usable email"})
		return
	}

	name := githubUser.Name
	if name == "" {
		name = githubUser.Login
	}

	authService := services.AuthService{}

	user, err := authService.HandleOAuthLogin(
		githubUser.Email,
		name,
		githubUser.Avatar,
		models.ProviderGitHub,
		strconv.FormatInt(githubUser.ID, 10),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Create persisted refresh session so refresh endpoint can validate + rotate tokens.
	session, err := authService.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("access_token", accessToken, int(utils.AccessTokenTTL.Seconds()), "/", "", true, true)
	c.SetCookie("refresh_token", session.RefreshToken, int(utils.RefreshTokenTTL.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"user": user})
}
