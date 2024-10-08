package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	S3Client *s3.Client
	Bucket   string
}

// NewS3Client initializes a new AWS S3 client
func NewS3Client(accessKeyId, accessKeySecret, endpoint, bucket string) (*S3Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Client{
		Bucket:   bucket,
		S3Client: client,
	}, nil
}

func (s *S3Client) GetPresignedURL(objectKey string, duration time.Duration) (string, error) {
	signer := s3.NewPresignClient(s.S3Client)

	request, err := signer.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})

	if err != nil {
		return "", err
	}

	return request.URL, err
}

func getContentTypeFromFileName(fileName string) string {
	if strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg") {
		return "image/jpeg"
	}

	if strings.HasSuffix(fileName, ".png") {
		return "image/png"
	}

	if strings.HasSuffix(fileName, ".gif") {
		return "image/gif"
	}

	if strings.HasSuffix(fileName, ".svg") {
		return "image/svg+xml"
	}

	if strings.HasSuffix(fileName, ".webp") {
		return "image/webp"
	}

	return "application/octet-stream"
}

func (s *S3Client) UploadFile(key string, file io.Reader) error {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(getContentTypeFromFileName(key)),
		Body:        file,
	})

	if err != nil {
		return err
	}

	return nil
}
