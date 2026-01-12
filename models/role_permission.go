package models

type RolePermission struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	RoleID     uint   `gorm:"not null;index:idx_role_permission,unique" json:"role_id"`
	Permission string `gorm:"size:150;not null;index:idx_role_permission,unique" json:"permission"`
}
