package main

import (
	"github.com/Kk120306/cvwo-2026/backend/internal/app"
)

// CompileDaemon --command="./backend"
// Main entry point for the application
func main() {
	// Create new application instance
	application := app.New()

	// Initialize application
	application.Initialize()

	// Run server
	application.Run()
}
