package eventconfig

import (
	"github.com/google/uuid"
)

type EventCertificateConfigResponse struct {
	ID                        uuid.UUID `json:"id"`
	EventID                   uuid.UUID `json:"event_id"`
	BaseCertificateStorageKey string    `json:"base_certificate_storage_key"`
	EventNamePosX             int32     `json:"event_name_pos_x"`
	EventNamePosY             int32     `json:"event_name_pos_y"`
	NamePosX                  int32     `json:"name_pos_x"`
	NamePosY                  int32     `json:"name_pos_y"`
	AcademicInstitutionPosX   int32     `json:"academic_institution_pos_x"`
	AcademicInstitutionPosY   int32     `json:"academic_institution_pos_y"`
	CreatedAt                 string    `json:"created_at"`
	UpdatedAt                 string    `json:"updated_at"`
}
