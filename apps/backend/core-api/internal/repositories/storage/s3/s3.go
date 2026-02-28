package s3

import (
	"apps/backend/common/utils"
	storage_datagateway "apps/backend/core-api/internal/datagateway/storage"
	"apps/backend/services/s3"
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
)

// S3Repository implements the S3DataGateway interface using the S3Service
type S3Repository struct {
	s3Service *s3.S3Service
}

// NewS3Repository creates a new S3Repository instance
func NewS3Repository(s3Service *s3.S3Service) *S3Repository {
	return &S3Repository{
		s3Service: s3Service,
	}
}

// GetStorageKey generates the storage key for a given entity
func (r *S3Repository) GetStorageKey(entityType storage_datagateway.StorageKeyType, entityID uuid.UUID, fileName string, fileExtension string) string {
	return r.s3Service.GetStorageKey(s3.StorageKeyType(entityType), entityID, fileName, fileExtension)
}

// GetS3UploadRequestObject prepares upload metadata from a file header
func (r *S3Repository) GetS3UploadRequestObject(entityType storage_datagateway.StorageKeyType, entityID uuid.UUID, fileHeader *multipart.FileHeader) (*storage_datagateway.S3UploadRequestObject, error) {
	// Use the underlying S3Service to create the request object
	// This maintains all the existing logic for file handling
	internalRequestObj, err := r.s3Service.GetS3UploadRequestObject(s3.StorageKeyType(entityType), entityID, fileHeader)
	if err != nil {
		return nil, err
	}

	// Generate storage key using the service's method
	extension := utils.GetFileExtension(fileHeader)
	storageKey := r.s3Service.GetStorageKey(s3.StorageKeyType(entityType), entityID, entityID.String(), extension)

	// Store the internal request object as opaque data
	// We'll use it directly in PutFile
	return &storage_datagateway.S3UploadRequestObject{
		EntityType:    entityType,
		EntityID:      entityID,
		StorageKey:    storageKey,
		FileHeader:    fileHeader,
		InternalData:  internalRequestObj, // Store the internal object for PutFile
	}, nil
}

// PutFile uploads a file to S3 and returns the storage key
func (r *S3Repository) PutFile(ctx context.Context, requestObject *storage_datagateway.S3UploadRequestObject) (string, error) {
	// Extract the internal S3Service request object
	internalRequestObj, ok := requestObject.InternalData.(*s3.S3UploadRequestObject)
	if !ok || internalRequestObj == nil {
		// Fallback: recreate the request object if not present
		var err error
		internalRequestObj, err = r.s3Service.GetS3UploadRequestObject(
			s3.StorageKeyType(requestObject.EntityType),
			requestObject.EntityID,
			requestObject.FileHeader,
		)
		if err != nil {
			return "", err
		}
	}

	return r.s3Service.PutFile(ctx, internalRequestObj)
}

// DeleteFile removes a file from S3 by its storage key
func (r *S3Repository) DeleteFile(ctx context.Context, key string) error {
	return r.s3Service.DeleteFile(ctx, key)
}

// GetPresignedURL generates a presigned URL for downloading a file
func (r *S3Repository) GetPresignedURL(ctx context.Context, key string) (string, error) {
	return r.s3Service.GetPresignedURL(ctx, key)
}

// GetFile retrieves a file from S3 as a readable stream
func (r *S3Repository) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return r.s3Service.GetFile(ctx, key)
}
