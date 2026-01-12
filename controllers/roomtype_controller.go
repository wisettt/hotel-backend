package controllers

import (
	"net/http"
	"hotel-backend/config"
	"hotel-backend/models"

	"github.com/gin-gonic/gin"
)

func GetRoomTypes(c *gin.Context) {
	var types []models.RoomType
	config.DB.Find(&types)
	c.JSON(http.StatusOK, types)
}

func CreateRoomType(c *gin.Context) {
	var rt models.RoomType
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&rt)
	c.JSON(http.StatusOK, rt)
}

func DeleteRoomType(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.RoomType{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Room type deleted"})
}
