package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Booking struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RoomID    *uint          `gorm:"column:room_id;index" json:"roomId,omitempty"`

	CustomerID       uint       `gorm:"index;column:customer_id" json:"customer_id"`
	ReferenceCode    string     `gorm:"column:reference_code;size:64" json:"reference_code,omitempty"`
	Status           string     `gorm:"column:status;size:64" json:"status,omitempty"`
	CheckIn          *time.Time `gorm:"column:check_in" json:"check_in,omitempty"`
	CheckOut         *time.Time `gorm:"column:check_out" json:"check_out,omitempty"`
	CheckInDate      *time.Time `gorm:"column:check_in_date" json:"check_in_date,omitempty"`
	CheckOutDate     *time.Time `gorm:"column:check_out_date" json:"check_out_date,omitempty"`
	Nights           int        `gorm:"column:nights" json:"nights,omitempty"`
	NumberOfGuests   int        `gorm:"column:number_of_guests" json:"number_of_guests,omitempty"`
	CheckinCompleted bool       `gorm:"column:checkin_completed;default:false" json:"checkinCompleted"`
	CheckedInAt      *time.Time `gorm:"column:checked_in_at" json:"checkedInAt,omitempty"`

	Adults   int `gorm:"column:adults;default:1" json:"adults"`
	Children int `gorm:"column:children;default:0" json:"children"`

	AccompanyingGuests datatypes.JSON `gorm:"column:accompanying_guests" json:"accompanyingGuests,omitempty"`

	Room     Room          `gorm:"foreignKey:RoomID;references:ID" json:"room,omitempty"`
	Customer Customer      `gorm:"foreignKey:CustomerID;references:ID" json:"customer,omitempty"`
	Rooms    []BookingRoom `gorm:"foreignKey:BookingID" json:"rooms"`
}
