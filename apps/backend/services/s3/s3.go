package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"time"

	"apps/backend/common/customerror"
	"apps/backend/common/log"
	"apps/backend/common/utils"
	config "apps/backend/core-api/config"
	s3Config "apps/backend/core-api/config/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type S3Service struct {
	s3Config   *s3Config.S3Config
	logger     *slog.Logger
	session    *session.Session
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3Service() (*S3Service, error) {
	cfg := config.LoadConfig()
	logger := log.LoadLogger()

	s3Cfg := &s3Config.S3Config{
		AccessKeyID:     cfg.S3.AccessKeyID,
		SecretAccessKey: cfg.S3.SecretAccessKey,
		BucketName:      cfg.S3.BucketName,
		Endpoint:        cfg.S3.Endpoint,
	}

	// Validate S3 configuration
	if !s3Cfg.IsValid() {
		logger.Error("S3 configuration is invalid or empty",
			"access_key_id", s3Cfg.AccessKeyID != "",
			"secret_access_key", s3Cfg.SecretAccessKey != "",
			"bucket_name", s3Cfg.BucketName,
			"endpoint", s3Cfg.Endpoint,
		)
		return nil, fmt.Errorf("S3 configuration is incomplete: please set S3_ACCESS_KEY_ID, S3_SECRET_ACCESS_KEY, S3_BUCKET_NAME, and S3_ENDPOINT in .env file")
	}

	// Create session following the Medium article pattern
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:                 aws.String("auto"),
			DisableParamValidation: aws.Bool(true),
			Endpoint:               aws.String(s3Cfg.Endpoint),
			Credentials:            credentials.NewStaticCredentials(s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey, ""),
			S3ForcePathStyle:       aws.Bool(true),
		},
	})
	if err != nil {
		logger.Error("Failed to create S3 session", "error", err)
		return nil, fmt.Errorf("failed to create S3 session: %w", err)
	}

	// Create S3 client and managers
	s3Client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &S3Service{
		s3Config:   s3Cfg,
		logger:     logger,
		session:    sess,
		s3Client:   s3Client,
		uploader:   uploader,
		downloader: downloader,
	}, nil
}

type (
	StorageKeyType        string
	S3UploadRequestObject struct {
		entityType    StorageKeyType
		entityID      uuid.UUID
		fileName      string
		contentType   string
		fileExtension string
		storageKey    string
		file          io.Reader
	}
)

const (
	StorageKeyTypeEventBanner      StorageKeyType = "event/%s/banner/%s%s"
	StorageKeyTypeEventIcon        StorageKeyType = "event/%s/icon/%s%s"
	StorageKeyTypeEventCertificate StorageKeyType = "event/%s/certificate/%s%s"
)

func (s *S3Service) GetStorageKey(entityType StorageKeyType, entityID uuid.UUID, fileName string, fileExtension string) string {
	return fmt.Sprintf(string(entityType), entityID, fileName, fileExtension)
}

func (s *S3Service) GetS3UploadRequestObject(entityType StorageKeyType, entityID uuid.UUID, fileHeader *multipart.FileHeader) (*S3UploadRequestObject, error) {
	file, err := utils.OpenFile(fileHeader)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	contentType := utils.GetFileContentType(fileHeader)
	extension := utils.GetFileExtension(fileHeader)
	fileName := entityID.String()
	storageKey := s.GetStorageKey(entityType, entityID, fileName, extension)

	return &S3UploadRequestObject{
		entityType:    entityType,
		entityID:      entityID,
		fileName:      fileName,
		fileExtension: extension,
		contentType:   contentType,
		storageKey:    storageKey,
		file:          file,
	}, nil
}

func (s *S3Service) PutFile(ctx context.Context, requestObject *S3UploadRequestObject) (string, error) {
	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(s.s3Config.BucketName),
		Key:         aws.String(requestObject.storageKey),
		Body:        requestObject.file,
		ContentType: aws.String(requestObject.contentType),
	})
	if err != nil {
		return "", customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return requestObject.storageKey, nil
}

func (s *S3Service) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.s3Config.BucketName),
		Key:    aws.String(key),
	}

	result, err := s.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return result.Body, nil
}

func (s *S3Service) DeleteFile(ctx context.Context, key string) error {
	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.s3Config.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return nil
}

func (s *S3Service) GetPresignedURL(ctx context.Context, key string) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.s3Config.BucketName),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(24 * time.Hour)
	if err != nil {
		return "", customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return urlStr, nil
}

func (s *S3Service) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.s3Config.BucketName),
		Prefix: aws.String(prefix),
	}

	result, err := s.s3Client.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	keys := make([]string, len(result.Contents))
	for i, obj := range result.Contents {
		keys[i] = *obj.Key
	}

	return keys, nil
}

func (s *S3Service) GetFileInfo(ctx context.Context, key string) (*s3.GetObjectOutput, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.s3Config.BucketName),
		Key:    aws.String(key),
	}

	result, err := s.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}

	return result, nil
}
