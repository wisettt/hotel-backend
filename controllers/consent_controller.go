package controllers

import (
	"fmt"
	
	"net/http"
	"strconv"
	"strings"
	"time"

	"hotel-backend/config"
	"hotel-backend/models"

	"github.com/gin-gonic/gin"
)

// -----------------------------
// Consents controller
// -----------------------------

// GET /api/consents
func GetConsents(c *gin.Context) {
	var consents []models.Consent
	if err := config.DB.Order("consent_id desc").Find(&consents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to load consents",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, consents)
}

// POST /api/consents
// POST /api/consents
func CreateConsent(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Slug        string `json:"slug" binding:"required"`
		Description string `json:"description"`
		Version     string `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid payload",
			"detail": err.Error(),
		})
		return
	}

	version := "1.0"
	if strings.TrimSpace(req.Version) != "" {
		version = strings.TrimSpace(req.Version)
	}

	// 🔑 1) หา consent เดิมก่อน
	var consent models.Consent
	err := config.DB.
		Where("slug = ? AND version = ?", req.Slug, version).
		First(&consent).Error

	if err != nil {
		// 🔑 2) ไม่เจอ → ค่อยสร้างใหม่
		consent = models.Consent{
			Title:       req.Title,
			Slug:        req.Slug,
			Description: req.Description,
			Version:     version,
		}

		if err := config.DB.Create(&consent).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to create consent",
				"detail": err.Error(),
			})
			return
		}
	}

	// ✅ คืน consent เดิมหรือใหม่ (แต่มีแถวเดียว)
	c.JSON(http.StatusOK, consent)
}

// ------------------------------------------------------------
// POST /api/consents/accept
// ------------------------------------------------------------
func AcceptConsent(c *gin.Context) {
    var req struct {
    GuestID   uint        `json:"guestId" binding:"required"`
    BookingID interface{} `json:"bookingId"` // ✅ ไม่ required
    ConsentID uint        `json:"consentId" binding:"required"`
    Action    string      `json:"action,omitempty"`
    Accepted  bool        `json:"accepted"`
}

    // 1️⃣ bind request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":  "invalid payload",
            "detail": err.Error(),
        })
        return
    }

    // 🔴 🔴 🔴 เพิ่มตรงนี้ 🔴 🔴 🔴
    var consent models.Consent
    if err := config.DB.First(&consent, req.ConsentID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid consentId",
        })
        return
    }
    // 🔴 🔴 🔴 จบตรงนี้ 🔴 🔴 🔴

    // 2️⃣ ตรวจ bookingId
   // bookingId อาจยังไม่รู้ → อนุญาตให้ nil
idPtr, _ := normalizeBookingIdentifier(req.BookingID)
    localGuestID := req.GuestID

action := req.Action
if strings.TrimSpace(action) == "" {
    action = "accepted"
}

status := "accepted"
if idPtr == nil {
    status = "pending"
}

cl := models.ConsentLog{
    ConsentID:  consent.ID,
    GuestID:    &localGuestID,
    BookingID:  idPtr,
    AcceptedAt: time.Now().UTC(),
    Status:     status,
    Action:     action,
}


    if err := config.DB.Create(&cl).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":  "db_error",
            "detail": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "ok":             true,
        "consent_log_id": cl.ID,
    })
}


// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------
func normalizeBookingIdentifier(raw interface{}) (*uint, string) {
	if raw == nil {
		return nil, ""
	}
	s := strings.TrimSpace(fmt.Sprint(raw))
	if s == "" {
		return nil, ""
	}
	if id64, err := strconv.ParseUint(s, 10, 64); err == nil {
		id := uint(id64)
		return &id, ""
	}
	return nil, s
}

// DELETE /api/consents/:id
func DeleteConsent(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}

	if err := config.DB.Delete(&models.Consent{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to delete consent",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "consent deleted"})
}
