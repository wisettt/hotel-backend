package models

import (
	"time"

	"gorm.io/gorm"
)

type BookingInfo struct {
	gorm.Model

	BookingID     uint       `gorm:"index" json:"bookingId"`
	Token         string     `gorm:"uniqueIndex;size:128" json:"token"`
	CheckinCode   string     `gorm:"uniqueIndex;size:16" json:"checkinCode"`
	Status        string     `gorm:"default:INITIATED" json:"status"`
	EmailStatus   string     `gorm:"default:PENDING" json:"emailStatus"`
	EmailError    string     `json:"emailError"`
	ExpiresAt     *time.Time `json:"expiresAt"`
	CodeExpiresAt *time.Time `json:"codeExpiresAt"`

	GuestEmail    string `json:"guestEmail"`
	GuestLastName string `json:"guestLastName"`
}
