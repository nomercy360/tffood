package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	S3Client *s3.Client
	Bucket   string
}

// NewS3Client initializes a new AWS S3 client
func NewS3Client(accessKeyId, accessKeySecret, endpoint, bucket string) (*Client, error) {
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

	return &Client{
		Bucket:   bucket,
		S3Client: client,
	}, nil
}

func (s *Client) GetPresignedURL(objectKey string, duration time.Duration) (string, error) {
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

func resolveContentType(fileName string) string {
	switch {
	case strings.HasSuffix(fileName, ".jpg"):
		return "image/jpeg"
	case strings.HasSuffix(fileName, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(fileName, ".png"):
		return "image/png"
	case strings.HasSuffix(fileName, ".gif"):
		return "image/gif"
	case strings.HasSuffix(fileName, ".webp"):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func (s *Client) UploadFile(file []byte, fileName string) (string, error) {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(resolveContentType(fileName)),
	})

	if err != nil {
		return "", err
	}

	return s.GetPresignedURL(fileName, time.Hour)
}
