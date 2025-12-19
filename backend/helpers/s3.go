package helpers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// config values
const (
	region     = "ap-southeast-2"
	bucketName = "direct-upload-s3-cvwo"
)

// GenerateUploadURL creates a presigned S3 PUT URL
func GenerateUploadURL() (string, error) {
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

	// Create presigner
	presigner := s3.NewPresignClient(client)

	// Generate random image name (16 bytes → hex)
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	imageName := hex.EncodeToString(randomBytes)

	// Presign PUT request
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageName),
	}, s3.WithPresignExpires(60*time.Second))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
