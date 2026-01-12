package storage

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

// StorageKeyType represents the type of storage path for different entities
type StorageKeyType string

const (
	StorageKeyTypeEventBanner      StorageKeyType = "event/%s/banner/%s%s"
	StorageKeyTypeEventIcon        StorageKeyType = "event/%s/icon/%s%s"
	StorageKeyTypeEventCertificate StorageKeyType = "event/%s/certificate/%s%s"
)

// S3UploadRequestObject encapsulates all information needed for an S3 upload
type S3UploadRequestObject struct {
	EntityType    StorageKeyType
	EntityID      uuid.UUID
	FileName      string
	ContentType   string
	FileExtension string
	StorageKey    string
	FileHeader    *multipart.FileHeader
	InternalData  interface{} // Opaque data used by repository implementation
}

// S3DataGateway defines the interface for S3 storage operations
type S3DataGateway interface {
	// GetS3UploadRequestObject prepares upload metadata from a file header
	GetS3UploadRequestObject(entityType StorageKeyType, entityID uuid.UUID, fileHeader *multipart.FileHeader) (*S3UploadRequestObject, error)

	// PutFile uploads a file to S3 and returns the storage key
	PutFile(ctx context.Context, requestObject *S3UploadRequestObject) (string, error)

	// DeleteFile removes a file from S3 by its storage key
	DeleteFile(ctx context.Context, key string) error

	// GetStorageKey generates the storage key for a given entity
	GetStorageKey(entityType StorageKeyType, entityID uuid.UUID, fileName string, fileExtension string) string

	// GetPresignedURL generates a presigned URL for downloading a file
	GetPresignedURL(ctx context.Context, key string) (string, error)

	// GetFile retrieves a file from S3 as a readable stream
	GetFile(ctx context.Context, key string) (interface{}, error) // Returns io.ReadCloser but defined as interface{} to avoid import
}
