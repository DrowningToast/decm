package certificate_share

import (
	"apps/backend/common/customerror"
	"context"

	"github.com/cockroachdb/errors"
)

// GetCertificateShareImage resolves a share by handle, enforces password protection,
// and returns the rendered PNG certificate image.
func (uc *CertificateShareUsecase) GetCertificateShareImage(ctx context.Context, handle string, password *string) ([]byte, error) {
	share, err := uc.CertificateShareDg.GetCertificateShareByHandle(ctx, handle)
	if err != nil {
		return nil, customerror.Parse(&customerror.ErrInternalServer, err)
	}
	if share == nil {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if !share.Active {
		return nil, customerror.Parse(&customerror.ErrNotFound, errors.New("certificate share not found"))
	}

	if err := uc.CheckSharePassword(share, password); err != nil {
		return nil, err
	}

	return uc.CertificateImageGenerator.GenerateCertificateImage(ctx, share.EventCertificateId)
}
