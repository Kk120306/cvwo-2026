package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// CommentRoutes set up the comment routes
func CommentRoutes(r *gin.Engine) {

	commentController := controllers.NewCommentController()

	commentRouter := r.Group("/comments") // Groups them under /comments
	{
		commentRouter.GET("/post/:postId", middleware.OptionalAuth, commentController.GetCommentsByPost)
		commentRouter.POST("/create/:postId", middleware.CheckAuth, commentController.CreateComment)
		commentRouter.DELETE("/delete/:id", middleware.CheckAuth, commentController.DeleteComment)
		commentRouter.PUT("/update/:id", middleware.CheckAuth, commentController.UpdateComment)
	}
}
