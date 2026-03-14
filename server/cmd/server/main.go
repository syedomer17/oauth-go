package main

import (
	"log"
	"oauth/internal/config"
	"oauth/internal/database"
	"oauth/internal/routes"
	"oauth/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.ConnectMongo()
	database.CreateIndexes(database.DB)
	oauth.InitGoogleOAuth()

	r := gin.Default()

	routes.AuthRoutes(r)

	log.Println("server is running on port",config.Config.AppPort)

	r.Run(":" + config.Config.AppPort)
}