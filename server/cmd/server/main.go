package main

import (
	"log"
	"oauth/internal/config"
	"oauth/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.ConnectMongo()
	database.CreateIndexes(database.DB)

	r := gin.Default()

	log.Println("server is running on port",config.Config.AppPort)

	r.Run(":" + config.Config.AppPort)
}