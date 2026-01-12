// models/customer.go
package models

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model

	// 🔥 แก้ไข: เปลี่ยน id ใน JSON tag ให้ Map ไปที่ string field
	// หรือให้เป็น field ที่ไม่ชนกับ GORM ID แต่เราต้องใช้ค่า string ชั่วคราวนี้
	FrontendTempID string `json:"frontendId,omitempty"`

	ID       uint   `gorm:"primaryKey"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	// ...
}

type CustomerBasic struct {
	ID       uint `gorm:"primaryKey"`
	FullName string
	Email    string
}
