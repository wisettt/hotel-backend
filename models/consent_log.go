package models

import (
    "time"

    "gorm.io/gorm"
)

type ConsentLog struct {
    ID           uint           `gorm:"primaryKey" json:"id"`
    BookingID    *uint          `gorm:"index" json:"booking_id"`
    BookingToken *string        `gorm:"type:varchar(255);index" json:"booking_token"`
    ConsentID    uint           `gorm:"index" json:"consent_id"`
    GuestID      *uint          `gorm:"index" json:"guest_id"`
    AcceptedAt   time.Time      `json:"accepted_at"`
    AcceptedBy   string         `json:"accepted_by"`
    Status       string         `gorm:"index" json:"status"`
    Action       string         `gorm:"index" json:"action"`

    CreatedAt    time.Time
    UpdatedAt    time.Time
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
