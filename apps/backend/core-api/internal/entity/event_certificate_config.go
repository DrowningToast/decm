package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventCertificateConfig struct {
	ID                              uuid.UUID `json:"id"`
	EventID                         uuid.UUID `json:"event_id"`
	BaseCertificateStorageKey       string    `json:"base_certificate_storage_key"`
	EventNamePosX                   float64   `json:"event_name_pos_x"`
	EventNamePosY                   float64   `json:"event_name_pos_y"`
	NamePosX                        float64   `json:"name_pos_x"`
	NamePosY                        float64   `json:"name_pos_y"`
	AcademicInstitutionPosX         *float64  `json:"academic_institution_pos_x"`
	AcademicInstitutionPosY         *float64  `json:"academic_institution_pos_y"`
	CertificateTitlePosX            *float64  `json:"certificate_title_pos_x"`
	CertificateTitlePosY            *float64  `json:"certificate_title_pos_y"`
	CertificateSubtitlePosX         *float64  `json:"certificate_subtitle_pos_x"`
	CertificateSubtitlePosY         *float64  `json:"certificate_subtitle_pos_y"`
	IsPublished                     bool      `json:"is_published"`
	EventNameFontFamilyID           *int32    `json:"event_name_font_family_id"`
	EventNameFontWeight             *int32    `json:"event_name_font_weight"`
	NameFontFamilyID                *int32    `json:"name_font_family_id"`
	NameFontWeight                  *int32    `json:"name_font_weight"`
	AcademicInstitutionFontFamilyID *int32    `json:"academic_institution_font_family_id"`
	AcademicInstitutionFontWeight   *int32    `json:"academic_institution_font_weight"`
	CertificateTitleFontFamilyID    *int32    `json:"certificate_title_font_family_id"`
	CertificateTitleFontWeight      *int32    `json:"certificate_title_font_weight"`
	CertificateSubtitleFontFamilyID *int32    `json:"certificate_subtitle_font_family_id"`
	CertificateSubtitleFontWeight   *int32    `json:"certificate_subtitle_font_weight"`
	CreatedAt                       time.Time `json:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at"`
}
