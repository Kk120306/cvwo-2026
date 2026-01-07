package main

import (
	"github.com/Kk120306/cvwo-2026/backend/internal/app"
	"log"
	"os"
)

// CompileDaemon --command="./backend"
// Main entry point for the application
func main() {
	log.SetOutput(os.Stdout)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    
    log.Println("========================================")
    log.Println("APPLICATION STARTING")
    log.Println("========================================")
	// Create new application instance
	application := app.New()

	log.Println("Calling app.Initialize()...")
	// Initialize application
	application.Initialize()
    log.Println("app.Initialize() COMPLETED")

	log.Println("Calling app.Run()...")
	// Run server
	application.Run()
	log.Println("app.Run() COMPLETED (should never reach here)")

}
