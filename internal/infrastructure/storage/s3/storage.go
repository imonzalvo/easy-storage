package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"easy-storage/internal/infrastructure/storage"
)

// S3Storage implements the Storage interface for AWS S3
type S3Storage struct {
	client *s3.Client
	bucket string
}

// NewS3Storage creates a new S3 storage implementation
func NewS3Storage(bucket, region, endpoint, accessKey, secretKey string) (*S3Storage, error) {
	// Create AWS SDK configuration
	opts := []func(*config.LoadOptions) error{}
	
	if accessKey != "" && secretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}
	
	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	
	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})
		opts = append(opts, config.WithEndpointResolverWithOptions(customResolver))
	}
	
	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	
	client := s3.NewFromConfig(awsCfg)
	
	return &S3Storage{
		client: client,
		bucket: bucket,
	}, nil
}

// PutObject uploads an object to S3
func (s *S3Storage) PutObject(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	
	if err != nil {
		return fmt.Errorf("failed to put object in S3: %w", err)
	}
	
	return nil
}

// GetObject retrieves an object from S3
func (s *S3Storage) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	
	return result.Body, nil
}

// DeleteObject removes an object from S3
func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}
	
	return nil
}

// GetObjectURL gets a presigned URL for direct access
func (s *S3Storage) GetObjectURL(ctx context.Context, key string, expires int64) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expires) * time.Second
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	
	return request.URL, nil
}
