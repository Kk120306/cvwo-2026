package helpers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)



// function invalidates Cloudfront cache for image
func InvalidateCloudFrontCache(imageName string) error {

	cloudfrontDistributionId:= os.Getenv("CLOUDFRONT_DISTRIBUTION_ID")
	ctx := context.Background()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region), // Uses the region from s3.go
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// From here on inspo is here 
	// https://www.trevorrobertsjr.com/blog/cloudfront-cache-invalidation-go/
	// https://www.youtube.com/watch?v=lZAGIy1e3JA&t=8s

	// Create CloudFront client
	cfClient := cloudfront.NewFromConfig(cfg)

	// Create invalidation and makes sure its unique by combining image name and time 
	callerReference := fmt.Sprintf("%s-%d", imageName, time.Now().Unix())

	// Invalidating the specific image 
	_, err = cfClient.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(cloudfrontDistributionId),
		InvalidationBatch: &types.InvalidationBatch{ // Define the invalidation batch
			CallerReference: aws.String(callerReference), // unique ID of the req
			Paths: &types.Paths{
				Quantity: aws.Int32(1), // one path
				Items: []string{
					"/" + imageName, // the path to invalidate 
				},
			},
		},
	})

	// If error occurs during invalidation
	if err != nil {
		return fmt.Errorf("CloudFront invalidation failed: %w", err)
	}

	return nil
}