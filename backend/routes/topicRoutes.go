package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// TopicRoutes sets up the topic routes
func TopicRoutes(r *gin.Engine) {

	topicController := controllers.NewTopicController()

	topicRouter := r.Group("/topics") // Groups them under /auth
	{
		topicRouter.GET("/", topicController.GetTopics)
		topicRouter.POST("/create", topicController.CreateTopic)
		// Only admin can update, delete and update topics
		topicRouter.DELETE("/delete/:slug", middleware.CheckAuth, middleware.CheckAdmin, topicController.DeleteTopic)
		topicRouter.PUT("/update/:slug", middleware.CheckAuth, middleware.CheckAdmin, topicController.UpdateTopic)
	}
}
