package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/database"
	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/Kk120306/cvwo-2026/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// retrives all topics avaliable sorted by creation date
func GetTopics(c *gin.Context) {
	// Create a slice for the topics
	var topics []string

	// query all the topics from database ordered - appends it to topic slice
	result := database.DB.Order("created_at desc").Find(&topics)

	// If there was a error in retriving from the database
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve topics",
		})
		return
	}

	// return the topics
	c.JSON(http.StatusOK, gin.H{
		"topics": topics,
	})
}

// Function that creates a topic and returns the topic object that was created
func CreateTopic(c *gin.Context) {
	// structure of the request body that we need
	var body struct {
		Name string
	}

	// if parsing the req to body fails, returns non nil, sends bad request
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	topic := models.Topic{
		Name: body.Name,
		Slug: helpers.GenerateSlug(body.Name),
	}
	// Creating the Topic
	result := database.DB.Create(&topic)

	// if there was an error during creation
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create topic",
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
func DeleteTopic(c *gin.Context) {
	// Get the slug of the topic
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Delete through Gorm - passing in the slug
	var topic models.Topic
	result := database.DB.Delete(&topic, "slug = ?", slug)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete topic",
		})
	}

	// When there were no changes to the database, means topic was not found
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic deleted successfully",
	})
}

// Function that updates a topic name
func UpdateTopic(c *gin.Context) {
	// Retrive the slug - name it cur as new topic name will change the slug
	curSlug := c.Param("slug")
	if curSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide a valid topic slug",
		})
		return
	}

	// Prase the new body that is passed - the new name user wants to set
	var body struct {
		Name string
	}

	// Must require empty validation in the frontend
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body, please provide valid details",
		})
		return
	}

	// Find the topic
	// https://gorm.io/docs/update.html
	var topic models.Topic
	result := database.DB.First(&topic, "slug = ?", curSlug)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Topic not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve topic",
		})
		return
	}

	// Updating the topic name and slug
	newSlug := helpers.GenerateSlug(body.Name)
	// Since topic is already found, it calls WHERE = loaded ID
	res := database.DB.Model(&topic).Updates(models.Topic{
		Name: body.Name,
		Slug: newSlug,
	})	
	
	// If there was an error during update
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update topic",
		})
		return
	}

	// Return the updated topic
	c.JSON(http.StatusOK, gin.H{
		"topic": topic,
	})

}
