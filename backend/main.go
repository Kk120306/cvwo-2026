package main

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/gin-gonic/gin"
	"github.com/Kk120306/cvwo-2026/backend/database"
)

// Loading enviornment variables & Connecting to database
func init() {
	helpers.LoadEnvVariables()
	database.ConnectToDb()
	database.PushDb()
}

// CompileDaemon --command="./backend"
// setting up gin project - https://gin-gonic.com/en/docs/quickstart/
func main() {
	router := gin.Default()

	// Testing route
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	router.Run() // listens on 0.0.0.0:8080 by default
}
