package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// PostRoutes sets up the post routes
func PostsRoutes(r *gin.Engine) {
	postsRouter := r.Group("/posts") // Groups them under /posts
	{
		postsRouter.GET("/all", middleware.OptionalAuth, controllers.GetAllPosts)
		postsRouter.GET("/topic/:slug", middleware.OptionalAuth, controllers.GetPostsByTopic)
		postsRouter.GET("/id/:id", middleware.OptionalAuth, controllers.GetPost)
		postsRouter.POST("/create/:slug", middleware.CheckAuth, controllers.CreatePost)
		postsRouter.DELETE("/delete/:id", middleware.CheckAuth, controllers.DeletePost)
		postsRouter.PUT("/update/:id", middleware.CheckAuth, controllers.UpdatePost)
		postsRouter.PATCH("/pin/:id", middleware.CheckAuth, middleware.CheckAdmin, controllers.TogglePinPost)
	}
}
