package models

import "time"

type Role struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Name        string            `gorm:"size:100;uniqueIndex" json:"name"`
	Description string            `gorm:"size:255" json:"description"`
	Permissions []RolePermission  `gorm:"foreignKey:RoleID" json:"permissions"`
	Members     []Admin           `gorm:"many2many:role_members;joinForeignKey:RoleID;JoinReferences:AdminID" json:"members"`
	CreatedAt   time.Time         `json:"created_at"`
}
