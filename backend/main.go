package main

import (
	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/Kk120306/cvwo-2026/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
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

	frontendURL := os.Getenv("FRONTEND_URL")

	// CORS configuration to allow requests from frontend
	// https://github.com/gin-contrib/cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setting all the routes 
	routes.AuthRoutes(router)
	routes.TopicRoutes(router)
	routes.PostsRoutes(router)
	routes.VoteRoutes(router)

	router.Run() // listens on 0.0.0.0:8080 by default
}
