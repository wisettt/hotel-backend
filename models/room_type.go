package models

import (
	"time"

	"gorm.io/gorm"
)

// RoomType Struct (นิยามเพียงแห่งเดียวในโปรเจกต์)
type RoomType struct {
	// 💡 Primary Key และ Timestamps
	ID uint `gorm:"primaryKey" json:"id"`

	TypeName    string `json:"typeName"`
	Description string `json:"description"`
	MaxGuests   uint   `json:"max_guests"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// One-To-Many Relation: RoomType -> Rooms
	// Rooms       []Room         `gorm:"foreignKey:RoomTypeID"`
}
