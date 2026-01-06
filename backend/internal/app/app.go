package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kk120306/cvwo-2026/backend/config"
	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// App struct represents the application
type App struct {
	Config *config.Config
	Router *gin.Engine
}

// New creates a new application instance
func New() *App {
	// Load configuration
	cfg := config.Load()

	// Validate configuration
	err := cfg.Validate()
	if err != nil {
		log.Fatal("Invalid configuration:", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	return &App{
		Config: cfg,
		Router: router,
	}
}

// Initialize sets up the application
func (a *App) Initialize() {
	// Connect to database
	database.ConnectToDb()

	// Push database migrations
	database.PushDb()

	// Uncomment for seeding
	// database.Seed()

	// Setup CORS middleware
	a.setupCORS()

	// Setup routes
	a.setupRoutes()
}

// configures CORS middleware
func (a *App) setupCORS() {
	// CORS configuration to allow requests from frontend
	// https://github.com/gin-contrib/cors
	a.Router.Use(cors.New(cors.Config{
		AllowOrigins:     a.Config.CORS.AllowedOrigins,
		AllowMethods:     a.Config.CORS.AllowedMethods,
		AllowHeaders:     a.Config.CORS.AllowedHeaders,
		ExposeHeaders:    a.Config.CORS.ExposeHeaders,
		AllowCredentials: a.Config.CORS.AllowCredentials,
		MaxAge:           a.Config.CORS.MaxAge,
	}))
}

// registers all application routes
func (a *App) setupRoutes() {
	// Setting all the routes
	routes.AuthRoutes(a.Router)
	routes.TopicRoutes(a.Router)
	routes.PostsRoutes(a.Router)
	routes.VoteRoutes(a.Router)
	routes.CommentRoutes(a.Router)
	routes.ImageRoutes(a.Router)
	routes.UserRoutes(a.Router)
}

// Run starts the application server
func (a *App) Run() {
	// Create server
	srv := &http.Server{
		Addr:    "0.0.0.0:" + a.Config.Server.Port,
		Handler: a.Router,
	}

	// Start server in a goroutine
	// ensures that server dosent block graceful shutdown handling
	go func() {
		log.Printf("Starting server on port %s in %s mode", a.Config.Server.Port, a.Config.Server.Env)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Gracefull shutdown
	// https://medium.com/@kittipat_1413/graceful-shutdown-in-golang-gin-a-complete-guide-130e3f075415

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
