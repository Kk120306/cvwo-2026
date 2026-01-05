package services

import (
	"errors"

	"github.com/Kk120306/cvwo-2026/backend/helpers"
)

// S3Service handles S3 operations business logic
type S3Service struct{}

// NewS3Service creates a new instance of S3Service
func NewS3Service() *S3Service {
	return &S3Service{}
}

// GenerateUploadURL generates a presigned URL for uploading to S3
func (s *S3Service) GenerateUploadURL() (string, error) {
	uploadURL, err := helpers.GenerateUploadURL()
	if err != nil {
		return "", errors.New("failed to generate S3 upload URL")
	}

	return uploadURL, nil
}

// DeleteImage deletes an image from S3 bucket
func (s *S3Service) DeleteImage(imageName string) error {
	// Validate image name
	if imageName == "" {
		return errors.New("image name is required")
	}

	// Deleting image from S3
	err := helpers.DeleteImage(imageName)
	if err != nil {
		return errors.New("failed to delete image from S3")
	}

	return nil
}