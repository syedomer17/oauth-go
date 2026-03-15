package routes

import (
	"oauth/internal/controller"
	"oauth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UploadRoutes(r *gin.Engine) {
	upload := r.Group("/upload")
	upload.POST("/", middleware.UploadMiddleware(), controller.UploadFile)
}