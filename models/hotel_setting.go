package models

import "time"

type HotelSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Email     string    `gorm:"size:150" json:"email"`
	Website   string    `gorm:"size:255" json:"website"`
	Logo      string    `gorm:"size:255" json:"logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
