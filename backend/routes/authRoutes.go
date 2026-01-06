package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up the authentication routes
func AuthRoutes(r *gin.Engine) {

	authController := controllers.NewAuthController()

	authRouter := r.Group("/auth") // Groups them under /auth
	{
		authRouter.POST("/signup", authController.Signup)
		authRouter.POST("/login", authController.Login)
		authRouter.POST("/logout", authController.Logout)
		authRouter.GET("/validate", middleware.CheckAuth, authController.Validate)
	}
}
