package main

import (
	"log"
	"oauth/internal/config"
	"oauth/internal/database"
	"oauth/internal/middleware"
	"oauth/internal/routes"
	"oauth/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.ConnectMongo()
	database.ConnectRedis()
	database.CreateIndexes(database.DB)
	oauth.InitGoogleOAuth()
	oauth.InitGithubOAuth()

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	routes.AuthRoutes(r)
	routes.UserRoutes(r)
	routes.UploadRoutes(r)

	log.Println("server is running on port", config.Config.AppPort)

	r.Run(":" + config.Config.AppPort)
}
