package services

import (
	"gorm.io/gorm"
	"hotel-backend/config"
	"hotel-backend/models"
)

// ConsentService manages consent templates.
type ConsentService struct {
	DB *gorm.DB
}

func NewConsentService(db *gorm.DB) *ConsentService {
	if db == nil {
		db = config.DB
	}
	return &ConsentService{DB: db}
}

func (s *ConsentService) Create(consent *models.Consent) error {
	if consent == nil {
		return gorm.ErrInvalidData
	}
	if consent.Version == "" {
		consent.Version = "1.0"
	}
	return s.DB.Create(consent).Error
}

func (s *ConsentService) List() ([]models.Consent, error) {
	var out []models.Consent
	// model uses consent_id as PK; order by consent_id desc
	err := s.DB.Order("consent_id desc").Find(&out).Error
	return out, err
}

func (s *ConsentService) GetByID(id uint) (models.Consent, error) {
	var c models.Consent
	err := s.DB.First(&c, "consent_id = ?", id).Error
	return c, err
}

func (s *ConsentService) Update(consent models.Consent) error {
	return s.DB.Model(&models.Consent{}).Where("consent_id = ?", consent.ID).Updates(consent).Error
}

func (s *ConsentService) Delete(id uint) error {
	return s.DB.Delete(&models.Consent{}, id).Error
}
