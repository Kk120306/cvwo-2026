package controllers

import (
	"net/http"

	"github.com/Kk120306/cvwo-2026/backend/helpers"
	"github.com/gin-gonic/gin"
)

// Function to get s3 upload Url
func GetS3UploadURL(c *gin.Context) {

	uploadURL, err := helpers.GenerateUploadURL()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate S3 upload URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
	})
}


// function to delete images from the s3 bucket 
func DeleteS3Image(c *gin.Context) {
	imageName := c.Param("imageName")

	if imageName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image name is required",
		})
		return
	}

	// Deleting image from S3
	err := helpers.DeleteImage(imageName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete image from S3",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}
