package middleware

import (
	"net/http"
	"oauth/internal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		maxSizeStr := config.Config.MaxFileSize

		maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)

		if err != nil {
			maxSize = 5 << 20 // 5MB default
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		file, header, err := c.Request.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "file required",
			})
			c.Abort()
			return
		}

		defer file.Close()

		if header.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "file too large",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}