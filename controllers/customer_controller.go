package controllers

import (
	"fmt"
	"hotel-backend/models"
	"hotel-backend/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	CustomerSvc *services.CustomerService
}

func NewCustomerController(svc *services.CustomerService) *CustomerController {
	return &CustomerController{CustomerSvc: svc}
}

// CreateCustomer (POST /api/customers) - T0.1
func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	// 💡 Note: ต้องใช้ 'var' เพื่อให้ GORM อัปเดตค่า ID เข้าไปใน Struct
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid customer payload: " + err.Error()})
		return
	}

	// 1. สร้าง Customer Record (Customer.ID จะถูกอัปเดตโดย GORM)
	if err := ctrl.CustomerSvc.Create(&customer); err != nil { // 💡 ส่ง Address (&customer)
		log.Printf("❌ DB ERROR during customer creation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Failed to create customer: %s", err.Error())})
		return
	}

	// 2. คืนค่า Customer Object พร้อม ID ที่ถูกสร้างโดย DB (HTTP 201 Created)
	// 💡 เนื่องจากเราไม่เห็นโค้ด models/customer.go ที่มีการเพิ่ม JSON tag 'id'
	//    เราจะสมมติว่า GORM/Gin จะแปลง ID (uint) เป็น 'id' ใน JSON อัตโนมัติ
	c.JSON(http.StatusCreated, customer)
}

// 💡 คุณสามารถเพิ่ม CRUD Methods อื่น ๆ ได้ที่นี่
