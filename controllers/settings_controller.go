package controllers

import (
	"errors"
	"net/http"

	"hotel-backend/config"
	"hotel-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type hotelSettingsPayload struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Logo    string `json:"logo"`
}

func GetHotelSettings(c *gin.Context) {
	var hotel models.HotelSetting
	if err := config.DB.First(&hotel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"hotel": models.HotelSetting{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotel": hotel})
}

func UpdateHotelSettings(c *gin.Context) {
	var payload hotelSettingsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hotel models.HotelSetting
	err := config.DB.First(&hotel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hotel = models.HotelSetting{
				Name:    payload.Name,
				Address: payload.Address,
				Phone:   payload.Phone,
				Email:   payload.Email,
				Website: payload.Website,
				Logo:    payload.Logo,
			}
			if err := config.DB.Create(&hotel).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"hotel": hotel})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hotel.Name = payload.Name
	hotel.Address = payload.Address
	hotel.Phone = payload.Phone
	hotel.Email = payload.Email
	hotel.Website = payload.Website
	hotel.Logo = payload.Logo

	if err := config.DB.Save(&hotel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotel": hotel})
}
