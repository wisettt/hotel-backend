package services

import (
	"hotel-backend/config"
	"hotel-backend/models"
)

type RoomService struct{}

func (s RoomService) Create(room models.Room) error {
	return config.DB.Create(&room).Error
}

func (s RoomService) GetAll() ([]models.Room, error) {
	var rooms []models.Room
	err := config.DB.Find(&rooms).Error
	return rooms, err
}

func (s RoomService) GetByID(id int) (models.Room, error) {
	var room models.Room
	// 💡 แก้ไข: ใช้ id ที่รับมาโดยตรง
	err := config.DB.First(&room, id).Error
	return room, err
}

func (s RoomService) Update(room models.Room) error {
	// 🔥 แก้ไข: เปลี่ยน room.RoomID เป็น room.ID
	return config.DB.Model(&models.Room{}).Where("id = ?", room.ID).Updates(room).Error
}

func (s RoomService) Delete(id int) error {
	return config.DB.Delete(&models.Room{}, id).Error
}
