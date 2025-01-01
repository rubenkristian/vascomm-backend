package utils

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rubenkristian/backend/configs"
)

type S3Option struct {
	S3Config *configs.S3Config
	s3Client minio.Client
	ctx      context.Context
}

func InitializeS3Client(s3Config *configs.S3Config) (*S3Option, error) {
	client, err := minio.New(s3Config.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecretKey, ""),
		Secure: s3Config.Ssl,
	})

	if err != nil {
		return nil, err
	}

	return &S3Option{
		S3Config: s3Config,
		s3Client: *client,
		ctx:      context.Background(),
	}, nil
}

func (s3Option *S3Option) UploadFile(file string, contentType string) (string, error) {
	s3Config := s3Option.S3Config
	info, err := s3Option.s3Client.FPutObject(s3Option.ctx, s3Config.BucketName, s3Config.FileLocation, file, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return "", err
	}

	return info.Location, nil
}

func (s3Option *S3Option) GetRegion() string {
	return s3Option.S3Config.Region
}

func (s3Option *S3Option) GetBucketName() string {
	return s3Option.S3Config.BucketName
}
