package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Response จาก Aigen
type AigenResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// --- 1. ID Card OCR (DoOCR) ---

func DoOCR(imageBase64 string) (map[string]interface{}, error) {
	// ... (ใช้โค้ดเดิมที่ดึง apiKey จาก os.Getenv) ...
	// Note: โค้ดนี้ใช้ os.Getenv("AIGEN_ENDPOINT") และ os.Getenv("AIGEN_API_KEY")
	// ...
	if _, err := base64.StdEncoding.DecodeString(imageBase64); err != nil {
		log.Println("[OCR] Base64 decode error:", err)
		return nil, fmt.Errorf("base64 invalid: %w", err)
	}

	endpoint := os.Getenv("AIGEN_ENDPOINT")
	apiKey := os.Getenv("AIGEN_API_KEY")

	payload := map[string]interface{}{
		"image": imageBase64,
		"model": "ocr-v1",
	}
	// ... (ส่วนสร้างและส่ง Request, Unmarshal, ตรวจสอบ Error) ...
	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("cannot build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-aigen-key", apiKey)

	// ... (ส่วนที่เหลือของ DoOCR) ... 
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}
	var ar AigenResponse
	if err := json.Unmarshal(bodyBytes, &ar); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}
	if ar.Status != "success" {
		return nil, fmt.Errorf("API status error: %s - %s", ar.Status, ar.Message)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(ar.Data, &arr); err == nil && len(arr) > 0 {
		return arr[0], nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(ar.Data, &obj); err == nil && len(obj) > 0 {
		return obj, nil
	}
	return nil, fmt.Errorf("no data returned from OCR: %s", string(ar.Data))
}

// --- 2. Passport OCR (DoPassportOCR) - โค้ดที่เพิ่มใหม่ ---

func DoPassportOCR(imageBase64 string) (map[string]interface{}, error) {

	if _, err := base64.StdEncoding.DecodeString(imageBase64); err != nil {
		log.Println("[Passport OCR] Base64 decode error:", err)
		return nil, fmt.Errorf("base64 invalid: %w", err)
	}

	// ดึงค่าจาก ENV
	endpoint := os.Getenv("AIGEN_ENDPOINT_PASSPORT")
	if endpoint == "" {
		endpoint = "https://api.aigen.online/aiscript/passport-ocr/v2"
	}
	apiKey := os.Getenv("AIGEN_API_KEY")

	payload := map[string]interface{}{
		"image": imageBase64,
		"model": "passport-ocr-v2",
	}

	b, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("cannot build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-aigen-key", apiKey)

	// ... (ส่วนสร้างและส่ง Request, Unmarshal, ตรวจสอบ Error เหมือนเดิม) ...
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ar AigenResponse
	if err := json.Unmarshal(bodyBytes, &ar); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	if ar.Status != "success" {
		return nil, fmt.Errorf("API status error: %s - %s", ar.Status, ar.Message)
	}

	var arr []map[string]interface{}
	if err := json.Unmarshal(ar.Data, &arr); err == nil && len(arr) > 0 {
		return arr[0], nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(ar.Data, &obj); err == nil && len(obj) > 0 {
		return obj, nil
	}

	return nil, fmt.Errorf("no data returned from Passport OCR: %s", string(ar.Data))
}
