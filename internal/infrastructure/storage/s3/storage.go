package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"easy-storage/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Provider implements the storage interface for S3-compatible storage
type S3Provider struct {
	client     *s3.Client
	bucketName string
}

// NewS3Provider creates a new S3 storage provider
func NewS3Provider(cfg *config.StorageConfig) (*S3Provider, error) {
	// Debug logging of configuration (remove in production)
	log.Printf("Initializing S3 provider with: Endpoint=%s, Bucket=%s, Region=%s",
		cfg.Endpoint, cfg.Bucket, cfg.Region)

	// Make sure endpoint has protocol
	endpoint := cfg.Endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
		log.Printf("Added http:// prefix to endpoint: %s", endpoint)
	}

	// Explicitly create static credentials
	staticCreds := credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"",
	)

	// Create a custom resolver for the endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: cfg.Region,
		}, nil
	})

	// Load configuration with explicit options
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(staticCreds),
	)
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return nil, err
	}

	// Create S3 client with path style option
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.ForcePathStyle
	})

	// Try a simple operation to verify connectivity
	_, err = client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		log.Printf("Error testing S3 connectivity: %v", err)
		return nil, fmt.Errorf("failed to connect to S3: %w", err)
	}

	log.Printf("Successfully connected to S3 endpoint")

	return &S3Provider{
		client:     client,
		bucketName: cfg.Bucket,
	}, nil
}

// Upload uploads a file to S3 storage
func (s *S3Provider) Upload(filename string, contentType string, file io.Reader) (string, error) {

	log.Printf("S3 Credentials - AccessKey: %s, SecretKey: %s", s.bucketName, context.TODO())
	// Generate a unique file path to avoid collisions
	ext := filepath.Ext(filename)
	uniquePath := time.Now().Format("2006/01/02/") + uuid.New().String() + ext

	// Upload file to S3
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(uniquePath),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return uniquePath, nil
}

// Download downloads a file from S3 storage
func (s *S3Provider) Download(path string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// Delete deletes a file from S3 storage
func (s *S3Provider) Delete(path string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	return err
}

// GetSignedURL generates a presigned URL for downloading a file
// expiryTime is the duration in seconds for which the URL will be valid
func (s *S3Provider) GetSignedURL(path string, expiryTime int64) (string, error) {
	// Create the presigned URL with an expiration time
	presignClient := s3.NewPresignClient(s.client)

	// Set expiration duration
	expires := time.Duration(expiryTime) * time.Second

	// Create the GetObject request
	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}
