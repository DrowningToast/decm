package certificate_share

import (
	"apps/backend/core-api/internal/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewCertificateShareViewModel(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	t.Run("returns nil when share is nil", func(t *testing.T) {
		result := NewCertificateShareViewModel(nil)
		assert.Nil(t, result)
	})

	t.Run("maps all fields from share without password", func(t *testing.T) {
		id := uuid.New()
		certID := uuid.New()
		share := &entity.CertificateShare{
			Id:                 id,
			EventCertificateId: certID,
			Active:             true,
			Handle:             "abc123handle",
			Password:           nil,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		vm := NewCertificateShareViewModel(share)

		assert.NotNil(t, vm)
		assert.Equal(t, id, vm.Id)
		assert.Equal(t, certID, vm.EventCertificateId)
		assert.True(t, vm.Active)
		assert.Equal(t, "abc123handle", vm.Handle)
		assert.False(t, vm.HasPassword)
		assert.Equal(t, now, vm.CreatedAt)
		assert.Equal(t, now, vm.UpdatedAt)
	})

	t.Run("sets HasPassword true when password is set", func(t *testing.T) {
		pw := "$argon2id$hashed"
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: uuid.New(),
			Password:           &pw,
		}

		vm := NewCertificateShareViewModel(share)

		assert.NotNil(t, vm)
		assert.True(t, vm.HasPassword)
	})

	t.Run("sets HasPassword false when password is nil", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: uuid.New(),
			Password:           nil,
		}

		vm := NewCertificateShareViewModel(share)

		assert.NotNil(t, vm)
		assert.False(t, vm.HasPassword)
	})

	t.Run("maps inactive share correctly", func(t *testing.T) {
		share := &entity.CertificateShare{
			Id:                 uuid.New(),
			EventCertificateId: uuid.New(),
			Active:             false,
			Handle:             "inactive-handle",
		}

		vm := NewCertificateShareViewModel(share)

		assert.NotNil(t, vm)
		assert.False(t, vm.Active)
		assert.Equal(t, "inactive-handle", vm.Handle)
	})
}
