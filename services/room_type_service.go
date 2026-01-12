package services

import (
	"hotel-backend/config"
	"hotel-backend/models"
)

type RoomTypeService struct{}

func (s RoomTypeService) Create(rt models.RoomType) error {
	return config.DB.Create(&rt).Error
}

func (s RoomTypeService) GetAll() ([]models.RoomType, error) {
	var types []models.RoomType
	err := config.DB.Find(&types).Error
	return types, err
}

func (s RoomTypeService) GetByID(id int) (models.RoomType, error) {
	var rt models.RoomType
	err := config.DB.First(&rt, id).Error
	return rt, err
}

func (s RoomTypeService) Update(rt models.RoomType) error {
	// 🔥 แก้ไข: เปลี่ยน rt.RoomTypeID เป็น rt.ID
	return config.DB.Model(&models.RoomType{}).Where("id = ?", rt.ID).Updates(rt).Error
}

func (s RoomTypeService) Delete(id int) error {
	return config.DB.Delete(&models.RoomType{}, id).Error
}