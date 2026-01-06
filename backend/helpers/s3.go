package helpers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// config values - inside url values so no need to hide in env
const (
	region     = "ap-southeast-2"
	bucketName = "direct-upload-s3-cvwo"
)

// GenerateUploadURL creates a presigned S3 PUT URL
// https://ronen-niv.medium.com/aws-s3-handling-presigned-urls-2718ab247d57
func GenerateUploadURL() (string, error) {
	// Create base context
	ctx := context.Background()

	// Load AWS config from environment variables
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)

	if err != nil {
		return "", err
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create presigner - wrapping client and generating signed url so s3 uploads possile from frontend.
	presigner := s3.NewPresignClient(client)

	// Generate random image name
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	imageName := hex.EncodeToString(randomBytes)
	// using crypt to secure random names

	// Create e url with embedded credential
	// make sure it expires after 60 seconds
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageName),
	}, s3.WithPresignExpires(60*time.Second))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

// function to delete image from s3 and also invalidate cloudfront cache
func DeleteImage(imageName string) error {
	// generate context
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return err
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Delete the object
	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageName),
	})

	// Invalidate cloudfront cache
	err = InvalidateCloudFrontCache(imageName)
	if err != nil {
		// Log the error but don't fail the request since S3 deletion succeeded
		log.Printf("CloudFront invalidation error: %v", err)
		// Cache can expire either ways so no need to really return error
	}

	return err
}
