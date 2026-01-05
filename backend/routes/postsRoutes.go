package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// PostRoutes sets up the post routes
func PostsRoutes(r *gin.Engine) {

	postController := controllers.NewPostController()

	postsRouter := r.Group("/posts") // Groups them under /posts
	{
		postsRouter.GET("/all", middleware.OptionalAuth, postController.GetAllPosts)
		postsRouter.GET("/topic/:slug", middleware.OptionalAuth, postController.GetPostsByTopic)
		postsRouter.GET("/id/:id", middleware.OptionalAuth, postController.GetPost)
		postsRouter.POST("/create/:slug", middleware.CheckAuth, postController.CreatePost)
		postsRouter.DELETE("/delete/:id", middleware.CheckAuth, postController.DeletePost)
		postsRouter.PUT("/update/:id", middleware.CheckAuth, postController.UpdatePost)
		postsRouter.PATCH("/pin/:id", middleware.CheckAuth, middleware.CheckAdmin, postController.TogglePinPost)
	}
}
