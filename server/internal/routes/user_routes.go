package routes

import (
	"oauth/internal/controller"
	"oauth/internal/middleware"

	"github.com/gin-gonic/gin"
)

// UserRoutes registers all user-facing endpoints under /api/user.
// Every route here is protected by the JWT cookie auth middleware.
func UserRoutes(r *gin.Engine) {
	user := r.Group("/api/user")
	user.Use(middleware.AuthRequired())

	// GET /api/user/me — returns the profile of the currently signed-in user.
	user.GET("/me", controller.GetMe)
}
