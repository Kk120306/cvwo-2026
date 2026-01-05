package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// TopicController handles HTTP requests for topics
type TopicController struct {
	topicService *services.TopicService
}

// NewTopicController creates a new instance of TopicController
func NewTopicController() *TopicController {
	return &TopicController{
		topicService: services.NewTopicService(),
	}
}

// retrieves all topics available sorted by creation date
func (tc *TopicController) GetTopics(c *gin.Context) {
	// Get topics through service layer
	topics, err := tc.topicService.GetAllTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return the topics
	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
	})
}

// Function that creates a topic and returns the topic object that was created
func (tc *TopicController) CreateTopic(c *gin.Context) {
	// structure of the request body that we need
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	// if parsing the req to body fails, returns non nil, sends bad request
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Create topic through service layer
	topic, err := tc.topicService.CreateTopic(services.CreateTopicInput{
		Name: body.Name,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Successfully created a topic so we just return a response here, return the topic
	c.JSON(http.StatusOK, gin.H{
		"topic": topic,
	})
}

// Function that deletes a topic
// https://gorm.io/docs/delete.html
func (tc *TopicController) DeleteTopic(c *gin.Context) {
	// Get the slug of the topic
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Delete through service layer
	err := tc.topicService.DeleteTopic(slug)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "topic not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
	})
}

// Function that updates a topic name
func (tc *TopicController) UpdateTopic(c *gin.Context) {
	// Retrieve the slug - name it cur as new topic name will change the slug
	curSlug := c.Param("slug")
	if curSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Parse the new body that is passed - the new name user wants to set
	var body struct {
		Name string `json:"name" binding:"required"`
	}

	// Must require empty validation in the frontend
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Find the topic through service layer
	topic, err := tc.topicService.FindTopicBySlug(curSlug)
	if err != nil {
		// Determine status code based on error type
		statusCode := http.StatusInternalServerError
		if err.Error() == "topic not found" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update topic through service layer
	err = tc.topicService.UpdateTopic(topic, services.UpdateTopicInput{
		Name: body.Name,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the updated topic
	c.JSON(http.StatusOK, gin.H{
		"topic": topic,
	})
}
