package services

import (
	"time"

	"gorm.io/gorm"
	"hotel-backend/config"
	"hotel-backend/models"
)

// ConsentLogService manages consent logs (accepted consents).
type ConsentLogService struct {
	DB *gorm.DB
}

func NewConsentLogService(db *gorm.DB) *ConsentLogService {
	if db == nil {
		db = config.DB
	}
	return &ConsentLogService{DB: db}
}

func (s *ConsentLogService) Log(cl *models.ConsentLog) error {
	if cl == nil {
		return gorm.ErrInvalidData
	}
	if cl.AcceptedAt.IsZero() {
		cl.AcceptedAt = time.Now().UTC()
	}
	if cl.Status == "" {
		if cl.BookingID != nil {
			cl.Status = "sent"
		} else {
			cl.Status = "pending"
		}
	}
	return s.DB.Create(cl).Error
}

func (s *ConsentLogService) LinkPendingByGuestIDs(bookingID uint, guestIDs []uint) (int64, error) {
	if len(guestIDs) == 0 {
		return 0, nil
	}
	now := time.Now().UTC()
	result := s.DB.Model(&models.ConsentLog{}).
		Where("booking_id IS NULL AND guest_id IN ?", guestIDs).
		Updates(map[string]interface{}{
			"booking_id": bookingID,
			"status":     "sent",
			"updated_at": now,
		})
	return result.RowsAffected, result.Error
}

func (s *ConsentLogService) LinkPendingByToken(token string, bookingID uint) (int64, error) {
	if token == "" {
		return 0, nil
	}
	now := time.Now().UTC()
	result := s.DB.Model(&models.ConsentLog{}).
		Where("booking_token = ? AND booking_id IS NULL", token).
		Updates(map[string]interface{}{
			"booking_id": bookingID,
			"status":     "sent",
			"updated_at": now,
		})
	return result.RowsAffected, result.Error
}
