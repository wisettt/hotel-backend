package models

import (
	

	"gorm.io/gorm"
)

type BookingRoom struct {
	gorm.Model
	BookingID uint  `gorm:"index;column:booking_id" json:"booking_id"`
	RoomID    uint  `gorm:"index;column:room_id" json:"room_id"`

	// optional details you used in service: nights, hours, status
	Nights int    `gorm:"column:nights;default:0" json:"nights,omitempty"`
	Hours  *int   `gorm:"column:hours" json:"hours,omitempty"`
	Status string `gorm:"column:status;size:64" json:"status,omitempty"`

	// timestamps already included via gorm.Model (CreatedAt, UpdatedAt, DeletedAt)
	// add convenience relation tags if needed:
	Booking Booking `gorm:"foreignKey:BookingID;references:ID" json:"booking,omitempty"`
	Room    Room    `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
}

// If you prefer not to include the Booking/Room relations to avoid recursive preloads,
// remove the Booking and Room fields above.
