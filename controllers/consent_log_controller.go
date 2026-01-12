package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotel-backend/config"
	"hotel-backend/models"
	"hotel-backend/utils"

	"github.com/gin-gonic/gin"
)

// GET /api/consent-logs
func GetConsentLogs(c *gin.Context) {
	var logs []models.ConsentLog

	if err := config.DB.Order("id desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load consent logs", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// POST /api/consent-logs
func CreateConsentLog(c *gin.Context) {
	var payload struct {
		BookingID       *uint   `json:"booking_id,omitempty"`
		BookingToken    *string `json:"booking_token,omitempty"`
		ConsentID        uint `json:"consentId" binding:"required"`
		GuestID         *uint   `json:"guest_id,omitempty"`
		Action          string  `json:"action,omitempty"`
		AcceptedAt      *string `json:"accepted_at,omitempty"` // accept string to parse multiple formats
		AcceptedBy      string  `json:"accepted_by,omitempty"`
		FaceImageBase64 *string `json:"face_image_base64,omitempty"` // optional base64 data URI or raw base64
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Printf("❌ Consent Log Binding Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}

	// parse accepted_at if provided (support RFC3339 or YYYY-MM-DD)
	acceptedAt := time.Now().UTC()
	if payload.AcceptedAt != nil && *payload.AcceptedAt != "" {
		if t, err := tryParseAcceptedAt(*payload.AcceptedAt); err == nil {
			acceptedAt = t
		} else {
			log.Printf("⚠️ could not parse accepted_at=%q: %v — using now", *payload.AcceptedAt, err)
		}
	}

	status := "pending"
	var bookingPtr *uint
	if payload.BookingID != nil && *payload.BookingID != 0 {
		bookingPtr = payload.BookingID
		status = "sent"
	}

	var tokenPtr *string
	if payload.BookingToken != nil && strings.TrimSpace(*payload.BookingToken) != "" {
		tmp := strings.TrimSpace(*payload.BookingToken)
		tokenPtr = &tmp
		// keep status pending if only token present
	}

	// try to save face image if provided
	var savedImagePath string
	if payload.FaceImageBase64 != nil && strings.TrimSpace(*payload.FaceImageBase64) != "" {
		// SaveBase64Image will accept either a data URI ("data:image/png;base64,...") or a raw base64 string
		if path, err := utils.SaveBase64Image(*payload.FaceImageBase64, "./uploads/faces"); err == nil {
			savedImagePath = path
			log.Printf("✅ saved face image: %s", path)
		} else {
			log.Printf("⚠️ could not save face image: %v", err)
			// do not block creation — just continue without image
		}
	}

	// Build entry
	entry := models.ConsentLog{
		BookingID:    bookingPtr,
		BookingToken: tokenPtr,
		ConsentID:    payload.ConsentID,
		AcceptedAt:   acceptedAt,
		AcceptedBy:   payload.AcceptedBy,
		Status:       status,
		Action:       payload.Action,
	}

	// set guest id if provided — models.ConsentLog.GuestID is *uint, so assign pointer directly
	if payload.GuestID != nil {
		entry.GuestID = payload.GuestID
	}

	// If we saved a face image and the Guest model has a FaceImagePath, update the Guest record.
	// This avoids adding a non-existent FaceImagePath field to ConsentLog.
	if savedImagePath != "" && payload.GuestID != nil {
		var g models.Guest
		if err := config.DB.First(&g, *payload.GuestID).Error; err == nil {
			g.FaceImagePath = savedImagePath
			if err := config.DB.Save(&g).Error; err != nil {
				log.Printf("⚠️ failed to save face image path to guest %d: %v", *payload.GuestID, err)
			}
		} else {
			log.Printf("⚠️ could not find guest %d to save face image path: %v", *payload.GuestID, err)
		}
	}

	if err := config.DB.Create(&entry).Error; err != nil {
		log.Printf("❌ DB Error creating consent_log: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create consent log", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}


// helper
func tryParseAcceptedAt(s string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		time.RFC1123,
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format")
}


// PATCH /api/consent-logs/attach-booking
// Accepts JSON: { bookingId: <number|string>, booking_id: <...>, guestIds?: [...], guestId?: <number> }
// PATCH /api/consent-logs/attach-booking
func AttachBookingToPending(c *gin.Context) {
	// read incoming JSON into a generic map so we can accept flexible keys/names
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		// if body is empty or invalid JSON, try to report useful message
		log.Printf("AttachBookingToPending: invalid payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload", "detail": err.Error()})
		return
	}

	// DEBUG: log payload keys (helps debugging missing bookingId)
	log.Printf("AttachBookingToPending payload keys: %+v", keysOfMap(req))

	// Accept bookingId from multiple keys or from query/header as fallback
	rawBooking, ok := req["bookingId"]
	if !ok {
		rawBooking, ok = req["booking_id"]
	}
	if !ok {
		// fallback: query param
		if qv := c.Query("bookingId"); qv != "" {
			rawBooking = qv
			ok = true
		} else if qv := c.Query("booking_id"); qv != "" {
			rawBooking = qv
			ok = true
		}
	}
	if !ok {
		// fallback: header
		if hv := c.GetHeader("X-Booking-Id"); hv != "" {
			rawBooking = hv
			ok = true
		}
	}

	if !ok {
		log.Printf("AttachBookingToPending missing bookingId - payload keys: %+v", keysOfMap(req))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "missing bookingId",
			"detail": "please provide bookingId (number or token) in request body (bookingId or booking_id), query, or X-Booking-Id header",
		})
		return
	}

	// parse booking identifier: either numeric id -> bookingID (uint) or token string -> bookingToken
	var bookingID *uint
	var bookingToken *string
	switch v := rawBooking.(type) {
	case float64:
		n := uint(v)
		bookingID = &n
	case int:
		n := uint(v)
		bookingID = &n
	case int64:
		n := uint(v)
		bookingID = &n
	case uint:
		n := v
		bookingID = &n
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bookingId empty string"})
			return
		}
		if id, err := strconv.ParseUint(s, 10, 64); err == nil {
			n := uint(id)
			bookingID = &n
		} else {
			bookingToken = &s
		}
	default:
		// try stringify fallback
		s := strings.TrimSpace(fmt.Sprint(v))
		if id, err := strconv.ParseUint(s, 10, 64); err == nil {
			n := uint(id)
			bookingID = &n
		} else {
			tmp := s
			if tmp == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bookingId empty or invalid"})
				return
			}
			bookingToken = &tmp
		}
	}

	// collect guest IDs: support guestIds array or single guestId (both camelCase and snake_case)
	var guestIDs []uint

	// 1) try guestIds
	if arr, ok := req["guestIds"]; ok {
		if cast, ok := arr.([]interface{}); ok {
			for _, it := range cast {
				switch vv := it.(type) {
				case float64:
					guestIDs = append(guestIDs, uint(vv))
				case string:
					if id, err := strconv.ParseUint(strings.TrimSpace(vv), 10, 64); err == nil {
						guestIDs = append(guestIDs, uint(id))
					}
				}
			}
		}
	}
	// 2) try guest_id (snake case) as array
	if len(guestIDs) == 0 {
		if arr, ok := req["guest_ids"]; ok {
			if cast, ok := arr.([]interface{}); ok {
				for _, it := range cast {
					switch vv := it.(type) {
					case float64:
						guestIDs = append(guestIDs, uint(vv))
					case string:
						if id, err := strconv.ParseUint(strings.TrimSpace(vv), 10, 64); err == nil {
							guestIDs = append(guestIDs, uint(id))
						}
					}
				}
			}
		}
	}
	// 3) try single guestId or guest_id
	if len(guestIDs) == 0 {
		var single interface{}
		if s, ok := req["guestId"]; ok {
			single = s
		} else if s, ok := req["guest_id"]; ok {
			single = s
		}
		if single != nil {
			switch vv := single.(type) {
			case float64:
				guestIDs = append(guestIDs, uint(vv))
			case string:
				if id, err := strconv.ParseUint(strings.TrimSpace(vv), 10, 64); err == nil {
					guestIDs = append(guestIDs, uint(id))
				}
			}
		}
	}

	if len(guestIDs) == 0 {
		log.Printf("AttachBookingToPending: no guestIds found in payload keys: %+v", keysOfMap(req))
		c.JSON(http.StatusBadRequest, gin.H{"error": "guestIds or guestId required"})
		return
	}

	// Debug: show parsed booking and guest ids
	log.Printf("AttachBookingToPending parsed bookingID=%v bookingToken=%v guestIDs=%v", bookingID, bookingToken, guestIDs)

	// Ensure there is at least one matching consent_log row to update (optional check to provide better feedback)
	var matchCount int64
	cond := config.DB.Model(&models.ConsentLog{}).Where("guest_id IN ?", guestIDs).Where("booking_id IS NULL")
	if err := cond.Count(&matchCount).Error; err != nil {
		log.Printf("AttachBookingToPending: failed to count matching consent_logs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error", "detail": err.Error()})
		return
	}
	if matchCount == 0 {
		log.Printf("AttachBookingToPending: no pending consent_logs found for guestIDs=%v", guestIDs)
		c.JSON(http.StatusOK, gin.H{"message": "no pending consent logs matched", "rows_affected": 0})
		return
	}

	// Build update map and perform update
	updateMap := map[string]interface{}{
		"status":     "sent",
		"updated_at": time.Now().UTC(),
	}
	if bookingID != nil {
		updateMap["booking_id"] = *bookingID
	} else if bookingToken != nil {
		updateMap["booking_token"] = *bookingToken
	}

	// Execute update
	res := cond.Updates(updateMap)
	if res.Error != nil {
		log.Printf("AttachBookingToPending update error: %v", res.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to attach booking", "detail": res.Error.Error()})
		return
	}

	log.Printf("AttachBookingToPending updated rows: %d", res.RowsAffected)
	c.JSON(http.StatusOK, gin.H{"message": "pending consent logs updated", "rows_affected": res.RowsAffected})
}

// DELETE /api/consent-logs/:id
func DeleteConsentLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := config.DB.Delete(&models.ConsentLog{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete consent log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consent log deleted"})
}

// helper to list map keys (for debugging logs)
func keysOfMap(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
