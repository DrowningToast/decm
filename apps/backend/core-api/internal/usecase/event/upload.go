package event

import (
	"context"
	"mime/multipart"

	"apps/backend/services/s3"

	"github.com/google/uuid"
)

func (uc *EventUsecase) UploadEventBanner(
	ctx context.Context,
	entityID uuid.UUID,
	bannerFile *multipart.FileHeader,
) (string, error) {
	requestObject, err := uc.S3Service.GetS3UploadRequestObject(s3.StorageKeyTypeEventBanner, entityID, bannerFile)
	if err != nil {
		return "", err
	}

	storageKey, err := uc.S3Service.PutFile(ctx, requestObject)
	if err != nil {
		return "", err
	}

	return storageKey, nil
}

func (uc *EventUsecase) UploadEventIcon(
	ctx context.Context,
	entityID uuid.UUID,
	iconFile *multipart.FileHeader,
) (string, error) {
	requestObject, err := uc.S3Service.GetS3UploadRequestObject(s3.StorageKeyTypeEventIcon, entityID, iconFile)
	if err != nil {
		return "", err
	}

	storageKey, err := uc.S3Service.PutFile(ctx, requestObject)
	if err != nil {
		return "", err
	}

	return storageKey, nil
}
