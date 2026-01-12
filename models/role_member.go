package models

type RoleMember struct {
	RoleID  uint `gorm:"primaryKey" json:"role_id"`
	AdminID uint `gorm:"primaryKey" json:"admin_id"`
}
