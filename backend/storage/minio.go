package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps the MinIO client
type MinIOClient struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

var Client *MinIOClient

// InitMinIO initializes the MinIO client with environment variables
func InitMinIO() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ROOT_USER")
	secretKey := os.Getenv("MINIO_ROOT_PASSWORD")
	bucketName := os.Getenv("MINIO_BUCKET")
	publicURL := os.Getenv("MINIO_PUBLIC_URL")

	// Default values for development
	if endpoint == "" {
		endpoint = "minio:9000"
	}
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	if bucketName == "" {
		bucketName = "uploads"
	}
	if publicURL == "" {
		publicURL = "http://localhost:9000"
	}

	// Initialize MinIO client
	ssl := os.Getenv("MINIO_SSL") == "true"
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		fmt.Printf("✅ Created MinIO bucket: %s\n", bucketName)
	}

	// Always set bucket policy to allow public read (runs every time to ensure policy is set)
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

	err = minioClient.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not set public policy: %v\n", err)
	} else {
		fmt.Printf("✅ Bucket '%s' set to public read access\n", bucketName)
	}

	Client = &MinIOClient{
		client:    minioClient,
		bucket:    bucketName,
		publicURL: publicURL,
	}

	fmt.Printf("✅ MinIO Storage initialized (endpoint: %s, bucket: %s)\n", endpoint, bucketName)
	return nil
}

// IsConfigured returns true if MinIO is properly configured
func IsConfigured() bool {
	return Client != nil && Client.client != nil
}

// UploadFile uploads a file to MinIO and returns (objectKey, publicURL, error)
func (m *MinIOClient) UploadFile(file *multipart.FileHeader, folder string) (string, string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFileName := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), randomString(8), ext)
	objectKey := fmt.Sprintf("%s/%s", folder, uniqueFileName)

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	_, err = m.client.PutObject(
		context.Background(),
		m.bucket,
		objectKey,
		src,
		file.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucket, objectKey)

	return objectKey, publicURL, nil
}

// UploadFromReader uploads from an io.Reader to MinIO
func (m *MinIOClient) UploadFromReader(reader io.Reader, objectKey string, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(
		context.Background(),
		m.bucket,
		objectKey,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	publicURL := fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucket, objectKey)
	return publicURL, nil
}

// DeleteFile deletes a file from MinIO
func (m *MinIOClient) DeleteFile(objectKey string) error {
	err := m.client.RemoveObject(
		context.Background(),
		m.bucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %w", err)
	}
	return nil
}

// GetPresignedURL generates a pre-signed URL for temporary access
func (m *MinIOClient) GetPresignedURL(objectKey string, expiration time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(
		context.Background(),
		m.bucket,
		objectKey,
		expiration,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

// GetPublicURL returns the public URL for an object
func (m *MinIOClient) GetPublicURL(objectKey string) string {
	return fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucket, objectKey)
}

// Helper function to generate random string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}
