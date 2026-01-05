package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoutes sets up the user routes
func UserRoutes(r *gin.Engine) {

	userController := controllers.NewUserController()

	userRouter := r.Group("/user") // Groups them under /user
	{
		userRouter.GET("/profile/:username", userController.GetUserProfile)
	}
}
