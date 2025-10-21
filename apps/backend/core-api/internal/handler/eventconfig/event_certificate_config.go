package eventconfig

import (
	"github.com/google/uuid"
)

type EventCertificateConfigResponse struct {
	ID                        uuid.UUID `json:"id"`
	EventID                   uuid.UUID `json:"event_id"`
	BaseCertificateStorageKey string    `json:"base_certificate_storage_key"`
	EventNamePosX             float64   `json:"event_name_pos_x"`
	EventNamePosY             float64   `json:"event_name_pos_y"`
	NamePosX                  float64   `json:"name_pos_x"`
	NamePosY                  float64   `json:"name_pos_y"`
	AcademicInstitutionPosX   *float64  `json:"academic_institution_pos_x"`
	AcademicInstitutionPosY   *float64  `json:"academic_institution_pos_y"`
	CreatedAt                 string    `json:"created_at"`
	UpdatedAt                 string    `json:"updated_at"`
}
