package models

import (
	"time"

	"gorm.io/gorm"
)

type Consent struct {
	// ใช้ ID ในโค้ด แต่แมปไปยัง column: consent_id ใน DB
	ID            uint           `gorm:"primaryKey;autoIncrement;column:consent_id" json:"id"`
	Slug          string         `json:"slug"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	EffectiveFrom *time.Time     `json:"effective_from"`
	Version       string         `json:"version"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
