package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
    "regexp"
	"hotel-backend/config"
	"hotel-backend/models"
	"hotel-backend/services"

	"github.com/gin-gonic/gin"
)

// --- Controller ---
type GuestController struct {
	GuestSvc *services.GuestService
}

// NewGuestController Constructor
func NewGuestController(svc *services.GuestService) *GuestController {
	return &GuestController{
		GuestSvc: svc,
	}
}

// Response struct
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ----------------------------------------------------------------------
// --- OCR ตรวจบัตรประชาชน ---
// ----------------------------------------------------------------------
func (c *GuestController) HandleIDCardVerification(ctx *gin.Context, apiKey string) {
	file, fileHeader, err := ctx.Request.FormFile("id_card_file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ไม่พบไฟล์ id_card_file"})
		return
	}
	defer file.Close()

	tempFilePath := fmt.Sprintf("%s/%s_%s", os.TempDir(), strconv.Itoa(os.Getpid()), fileHeader.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ไม่สามารถสร้าง temp file ได้"})
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFilePath)

	_, _ = io.Copy(tempFile, file)

	b, _ := os.ReadFile(tempFilePath)
	imageBase64 := base64.StdEncoding.EncodeToString(b)

	result, err := services.DoOCR(imageBase64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "OCR สำเร็จ",
		Data:    result,
	})
}

// ----------------------------------------------------------------------
// --- Passport OCR ---
// ----------------------------------------------------------------------
func (c *GuestController) HandlePassportVerification(ctx *gin.Context, apiKey string) {
	file, fileHeader, err := ctx.Request.FormFile("passport_file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "ไม่พบไฟล์ passport_file"})
		return
	}
	defer file.Close()

	tempFilePath := fmt.Sprintf("%s/%s_%s", os.TempDir(), strconv.Itoa(os.Getpid()), fileHeader.Filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ไม่สามารถสร้าง temp file ได้"})
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFilePath)

	_, _ = io.Copy(tempFile, file)

	b, _ := os.ReadFile(tempFilePath)
	imageBase64 := base64.StdEncoding.EncodeToString(b)

	result, err := services.DoPassportOCR(imageBase64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Passport OCR สำเร็จ",
		Data:    result,
	})
}

// ----------------------------------------------------------------------
// ✅ NEW: Get Guests by Booking ID (ชัวร์ว่าไม่ปน)
// GET /api/bookings/:id/guests
// ----------------------------------------------------------------------
func (c *GuestController) GetGuestsByBookingID(ctx *gin.Context) {
    log.Println("✅ HIT GetGuestsByBookingID id=", ctx.Param("id"))

    idStr := strings.TrimSpace(ctx.Param("id"))
    if idStr == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "error.missingBookingId",
        })
        return
    }

    bookingID64, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil || bookingID64 == 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "error.invalidBookingId",
        })
        return
    }
    bookingID := uint(bookingID64)

    var guests []models.Guest

    // ✅ FILTER booking_id = bookingID เท่านั้น
    if err := config.DB.
        Where("booking_id = ?", bookingID).
        Order("is_main_guest DESC, id ASC").
        Find(&guests).Error; err != nil {

        log.Printf("[GetGuestsByBookingID] booking_id=%d err=%v", bookingID, err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": "Failed to fetch guests",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   guests,
    })
}

// ----------------------------------------------------------------------
// ⚠️ Get Guests -> บังคับกรองด้วย bookingId เท่านั้น
// GET /api/guests?bookingId=123
// ----------------------------------------------------------------------
func (c *GuestController) GetGuests(ctx *gin.Context) {
	log.Println("✅ HIT GetGuests bookingId=", ctx.Query("bookingId"))

	q := strings.TrimSpace(ctx.Query("bookingId"))
	if q == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "error.missingBookingId",
		})
		return
	}

	bookingID64, err := strconv.ParseUint(q, 10, 64)
	if err != nil || bookingID64 == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "error.invalidBookingId",
		})
		return
	}
	bookingID := uint(bookingID64)

	var guests []models.Guest
	if err := config.DB.
		Where("booking_id = ?", bookingID).
		Order("is_main_guest DESC, id ASC").
		Find(&guests).Error; err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch guests",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   guests,
	})
}

// ----------------------------------------------------------------------
// GET /api/guests/all
// ----------------------------------------------------------------------
func (c *GuestController) GetAllGuests(ctx *gin.Context) {
	log.Println("✅ HIT GetAllGuests")

	// แนะนำให้ใช้ service (คุณมี GetAll() แล้ว) เพื่อ preload / แต่งข้อมูลได้
	guests, err := c.GuestSvc.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   guests,
	})
}

func (c *GuestController) GetGuestByID(ctx *gin.Context) {
    log.Println("✅ HIT GetGuestByID id=", ctx.Param("id"))

    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid guest ID"})
        return
    }

    guest, err := c.GuestSvc.GetByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Guest not found"})
        return
    }

    ctx.JSON(http.StatusOK, guest)
}

func (c *GuestController) UpdateGuest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid guest ID"})
		return
	}

	var payload models.Guest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	payload.ID = uint(id)

	if err := c.GuestSvc.Update(&payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, payload)
}

// ----------------------------------------------------------------------
// --- CreateGuest (STRICT: bookingId only) ---
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// --- CreateGuest (รองรับ camelCase + snake_case) ---
// POST /api/guests
// ----------------------------------------------------------------------
func (c *GuestController) CreateGuest(ctx *gin.Context) {
	var payload map[string]interface{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("➡️ RAW CreateGuest payload: %+v", payload)

	// ---------------- helpers ----------------
	getString := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := payload[k]; ok && v != nil {
				switch vv := v.(type) {
				case string:
					s := strings.TrimSpace(vv)
					if s != "" {
						return s
					}
				case fmt.Stringer:
					s := strings.TrimSpace(vv.String())
					if s != "" {
						return s
					}
				}
			}
		}
		return ""
	}

	getBool := func(keys ...string) (bool, bool) {
		for _, k := range keys {
			if v, ok := payload[k]; ok && v != nil {
				switch vv := v.(type) {
				case bool:
					return vv, true
				case float64:
					return vv != 0, true
				case int:
					return vv != 0, true
				case string:
					s := strings.TrimSpace(strings.ToLower(vv))
					if s == "true" || s == "1" || s == "yes" {
						return true, true
					}
					if s == "false" || s == "0" || s == "no" {
						return false, true
					}
				}
			}
		}
		return false, false
	}

	getUintPtr := func(keys ...string) (*uint, bool) {
		for _, k := range keys {
			if v, ok := payload[k]; ok && v != nil {
				switch vv := v.(type) {
				case float64:
					n := uint(vv)
					if n > 0 {
						return &n, true
					}
				case int:
					n := uint(vv)
					if n > 0 {
						return &n, true
					}
				case string:
					s := strings.TrimSpace(vv)
					if s == "" {
						continue
					}
					if id, err := strconv.ParseUint(s, 10, 64); err == nil && id > 0 {
						n := uint(id)
						return &n, true
					}
				}
			}
		}
		return nil, false
	}

	parseDOB := func(v string) *time.Time {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil
		}

		// รองรับ "2006-01-02"
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return &t
		}

		// รองรับ RFC3339 เช่น "2004-07-23T07:00:00+07:00"
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			return &tt
		}

		log.Println("⚠️ invalid dateOfBirth/date_of_birth:", v)
		return nil
	}

	// ---------------- bookingId (required) ----------------
	bookingID, ok := getUintPtr("bookingId", "booking_id")
	if !ok || bookingID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "error.missingBookingId",
		})
		return
	}

	// ---------------- map payload -> model ----------------
	var g models.Guest
	g.BookingID = bookingID

	// full name
	fullName := getString("fullName", "full_name", "name")
	if fullName != "" {
		g.FullName = fullName
	}

	// is main guest
	if b, ok := getBool("isMainGuest", "is_main_guest", "mainGuest", "main_guest"); ok {
		g.IsMainGuest = b
	}

	// date of birth
	dob := getString("dateOfBirth", "date_of_birth")
	if dob != "" {
		g.DateOfBirth = parseDOB(dob)
	}

	// basic info
	g.Gender = getString("gender")
	g.Nationality = getString("nationality")

	// address
	addr := getString("currentAddress", "current_address")
	if addr != "" {
		g.CurrentAddress = addr
	}

	// document type / number
	// React ส่ง documentType: "ID_CARD" | "PASSPORT"
	docType := getString("documentType", "idType", "id_type")
	if docType != "" {
		g.IDType = docType
	}

	docNo := getString("documentNumber", "idNumber", "id_number")
	if docNo != "" {
		g.IDNumber = docNo
	}

	issued := getString("idIssuedCountry", "id_issued_country")
	if issued != "" {
		g.IDIssuedCountry = issued
	}

	// images path (ถ้าส่งเป็น path มาอยู่แล้ว)
	g.FaceImagePath = getString("faceImagePath", "face_image_path")
	g.DocumentImagePath = getString("documentImagePath", "document_image_path")

	// base64 images (React ส่ง faceImageBase64/documentImageBase64)
	faceB64 := getString("faceImageBase64", "face_image_base64")
	if faceB64 != "" {
		path, err := services.SaveBase64Image(faceB64, "faces")
		if err != nil {
			log.Println("❌ save face image failed:", err)
		} else {
			g.FaceImagePath = path
		}
	}

	docB64 := getString("documentImageBase64", "document_image_base64")
	if docB64 != "" {
		path, err := services.SaveBase64Image(docB64, "documents")
		if err != nil {
			log.Println("❌ save document image failed:", err)
		} else {
			g.DocumentImagePath = path
		}
	}

	log.Printf("➡️ CreateGuest mapped model: %+v", g)

	// ---------------- save ----------------
	if err := c.GuestSvc.Create(&g); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   g,
	})
}

// ----------------------------------------------------------------------
// --- Delete Guest ---
// ----------------------------------------------------------------------
func (c *GuestController) DeleteGuest(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, gin.H{
		"status":  "error",
		"message": "guest deletion is disabled",
	})
}


// ฟังก์ชันตรวจสอบอีเมลที่ถูกต้อง
func isValidEmail(email string) bool {
    // Regular Expression สำหรับตรวจสอบอีเมล
    re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}
