package services

import (
	"hotel-backend/config"
	"hotel-backend/models"
)

type AdminService struct{}

func (s AdminService) Create(admin models.Admin) error {
	return config.DB.Create(&admin).Error
}

func (s AdminService) GetAll() ([]models.Admin, error) {
	var admins []models.Admin
	err := config.DB.Find(&admins).Error
	return admins, err
}

func (s AdminService) GetByID(id int) (models.Admin, error) {
	var admin models.Admin
	err := config.DB.First(&admin, id).Error
	return admin, err
}

type Admin struct {
    AdminID uint `gorm:"primaryKey" json:"admin_id"` // ⬅️ ชื่อฟิลด์คือ AdminID
    // ...
}

func (s AdminService) Delete(id int) error {
	return config.DB.Delete(&models.Admin{}, id).Error
}