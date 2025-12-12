package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// TopicRoutes sets up the topic routes
func TopicRoutes(r *gin.Engine) {
	topicRouter := r.Group("/topics") // Groups them under /auth
	{
		topicRouter.GET("/", controllers.GetTopics) 
		// Only admin can create, delete and update topics
		topicRouter.POST("/create", middleware.CheckAuth, middleware.CheckAdmin, controllers.CreateTopic)
		topicRouter.DELETE("/delete/:slug", middleware.CheckAuth, middleware.CheckAdmin, controllers.DeleteTopic)
		topicRouter.PUT("/update/:slug", middleware.CheckAuth, middleware.CheckAdmin, controllers.UpdateTopic)
	}
}
