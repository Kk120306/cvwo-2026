package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/services"
	"github.com/gin-gonic/gin"
)

// S3Controller handles HTTP requests for S3 operations
type S3Controller struct {
	s3Service *services.S3Service
}

// NewS3Controller creates a new instance of S3Controller
func NewS3Controller() *S3Controller {
	return &S3Controller{
		s3Service: services.NewS3Service(),
	}
}

// Function to get s3 upload URL
func (sc *S3Controller) GetS3UploadURL(c *gin.Context) {
	// Generate upload URL through service layer
	uploadURL, err := sc.s3Service.GenerateUploadURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uploadUrl": uploadURL,
	})
}

// function to delete images from the s3 bucket
func (sc *S3Controller) DeleteS3Image(c *gin.Context) {
	imageName := c.Param("imageName")

	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image name is required",
		})
		return
	}

	// Delete image through service layer
	err := sc.s3Service.DeleteImage(imageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
