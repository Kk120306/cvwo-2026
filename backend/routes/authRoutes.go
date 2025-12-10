package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// AuthRoutes sets up the authentication routes
func AuthRoutes(r *gin.Engine) {
	authRouter := r.Group("/auth") // Groups them under /auth
	{
		authRouter.POST("/signup", controllers.Signup)
		authRouter.POST("/login", controllers.Login)
		authRouter.POST("/logout", controllers.Logout)
		authRouter.GET("/validate", middleware.CheckAuth, controllers.Validate)
	}
}
