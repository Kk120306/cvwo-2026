package routes

import (
	"github.com/Kk120306/cvwo-2026/backend/controllers"
	"github.com/Kk120306/cvwo-2026/backend/middleware"
	"github.com/gin-gonic/gin"
)

// VoteRoutes sets up the vote routes
func VoteRoutes(r *gin.Engine) {
	voteRouter := r.Group("/vote") // Groups them under /vote
	{
		voteRouter.POST("/", middleware.CheckAuth, controllers.CreateOrUpdateVote)
		// id is the content Id and type is either "post" or "comment"
		voteRouter.GET("/count/:id/:type", controllers.GetVotesCount)
	}
}
