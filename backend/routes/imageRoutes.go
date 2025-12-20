package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// ImageRoutes sets up the post routes
func ImageRoutes(r *gin.Engine) {
	imageRouter := r.Group("/images") // Groups them under /images
	{
		imageRouter.GET("/s3Url", middleware.CheckAuth, controllers.GetS3UploadURL)
		imageRouter.DELETE("/delete/:imageName", middleware.CheckAuth, controllers.DeleteS3Image)
	}
}
