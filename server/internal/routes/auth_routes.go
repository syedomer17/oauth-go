package routes

import (
	"oauth/internal/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")

	auth.GET("/google/login", controller.GoogleLogin)
	auth.GET("/google/callback", controller.GoogleCallback)
	auth.GET("/github/login", controller.GithubLogin)
	auth.GET("/github/callback", controller.GithubCallback)
	auth.POST("/refresh", controller.RefreshToken)
	auth.POST("/logout", controller.Logout)
}
