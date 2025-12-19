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
