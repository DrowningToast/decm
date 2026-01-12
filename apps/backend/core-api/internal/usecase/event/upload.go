package event

import (
	storage_datagateway "apps/backend/core-api/internal/datagateway/storage"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

func (uc *EventUsecase) uploadEventBannerImpl(
	ctx context.Context,
	entityID uuid.UUID,
	bannerFile *multipart.FileHeader,
) (string, error) {
	requestObject, err := uc.S3DataGateway.GetS3UploadRequestObject(storage_datagateway.StorageKeyTypeEventBanner, entityID, bannerFile)
	if err != nil {
		return "", err
	}

	storageKey, err := uc.S3DataGateway.PutFile(ctx, requestObject)
	if err != nil {
		return "", err
	}

	return storageKey, nil
}

func (uc *EventUsecase) uploadEventIconImpl(
	ctx context.Context,
	entityID uuid.UUID,
	iconFile *multipart.FileHeader,
) (string, error) {
	requestObject, err := uc.S3DataGateway.GetS3UploadRequestObject(storage_datagateway.StorageKeyTypeEventIcon, entityID, iconFile)
	if err != nil {
		return "", err
	}

	storageKey, err := uc.S3DataGateway.PutFile(ctx, requestObject)
	if err != nil {
		return "", err
	}

	return storageKey, nil
}
