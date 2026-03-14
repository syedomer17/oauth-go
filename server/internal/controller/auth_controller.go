package controller

import (
	"net/http"
	"oauth/internal/models"
	"oauth/internal/services"
	"oauth/internal/utils"
	"oauth/pkg/oauth"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func GoogleLogin(c *gin.Context) {
	state := utils.GenerateState()

	url := oauth.GoogleOAuthConfig.AuthCodeURL(state)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := oauth.GoogleOAuthConfig.Exchange(c,code)

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
		ID string `json:"id"`
		Email string `json:"email"`
		Name string `json:"name"`
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

	jwtToken, err := utils.GenerateJWT(user.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": jwtToken,
		"user": user,
	})

}