package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FullName  string         `gorm:"size:255" json:"full_name"`
	Username  string         `gorm:"uniqueIndex;size:150" json:"username"`
	Password  string         `gorm:"size:255" json:"-"` // store hashed password, never return in JSON
	ResetToken        *string    `gorm:"size:128;index" json:"-"`
	ResetTokenExpires *time.Time `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
