package handlers

import (
	"net/http"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateLuggageRequest 行李寄存请求结构体
type CreateLuggageRequest struct {
	GuestName    string `json:"guest_name" binding:"required"` // 客人姓名
	ContactPhone string `json:"contact_phone"`                 // 联系电话
	ContactEmail string `json:"contact_email"`                 // 联系邮箱
	Description  string `json:"description"`                   // 行李描述
	Quantity     int    `json:"quantity"`                      // 行李数量
	SpecialNotes string `json:"special_notes"`                 // 特殊备注
	StoreroomID  int64  `json:"storeroom_id" binding:"required"`// 寄存室ID
	StoredBy     int64  `json:"stored_by" binding:"required"`   // 操作员ID
	QRCodeURL    string `json:"qr_code_url"`                   // 二维码URL（可选）
}

// CreateLuggage 处理行李寄存请求
func CreateLuggage(c *gin.Context) {
	var req CreateLuggageRequest
	// 解析并校验 JSON 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	item, err := services.CreateLuggage(services.CreateLuggageRequest{
		GuestName:    req.GuestName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Description:  req.Description,
		Quantity:     req.Quantity,
		SpecialNotes: req.SpecialNotes,
		StoreroomID:  req.StoreroomID,
		StoredBy:     req.StoredBy,
		QRCodeURL:    req.QRCodeURL,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "create luggage success",
		"luggage_id":     item.ID,
		"retrieval_code": item.RetrievalCode,
	})
}
