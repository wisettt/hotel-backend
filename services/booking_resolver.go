package services

import (
    "errors"
    "fmt"
    "strings"

    "hotel-backend/config"
)

func ResolveBookingIDFromToken(token string) (uint, error) {
    token = strings.TrimSpace(token)
    if token == "" {
        return 0, errors.New("empty token")
    }

    // ปรับ query ให้ตรงกับ schema ของคุณ: ตัวอย่างนี้สมมติมี table booking_infos หรือ bookings with token column
    var out struct {
        ID uint `gorm:"column:id"`
    }

    // ตัวอย่าง: ถ้าคุณเก็บ token ใน table booking_infos column checkin_token
    // SELECT booking_info_id or id depending on schema; ปรับ SQL ให้ตรง
    if err := config.DB.Raw("SELECT id FROM booking_infos WHERE checkin_token = ? LIMIT 1", token).Scan(&out).Error; err != nil {
        return 0, fmt.Errorf("db query failed: %w", err)
    }

    if out.ID == 0 {
        return 0, fmt.Errorf("not found")
    }
    return out.ID, nil
}
