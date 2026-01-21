package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateLuggageRequest 行李寄存请求结构体
type CreateLuggageRequest struct {
	GuestName    string `json:"guest_name" binding:"required"`   // 客人姓名
	ContactPhone string `json:"contact_phone"`                   // 联系电话
	ContactEmail string `json:"contact_email"`                   // 联系邮箱
	Description  string `json:"description"`                     // 行李描述
	Quantity     int    `json:"quantity"`                        // 行李数量
	SpecialNotes string `json:"special_notes"`                   // 特殊备注
	StoreroomID  int64  `json:"storeroom_id" binding:"required"` // 寄存室ID
	StoredBy     int64  `json:"stored_by" binding:"required"`    // 操作员ID
	QRCodeURL    string `json:"qr_code_url"`                     // 二维码URL（可选）
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

// QueryLuggageByUserInfo 按用户信息查询寄存记录
// GET /storage/search?guest_name=...&contact_phone=...
func QueryLuggageByUserInfo(c *gin.Context) {
	guestName := c.Query("guest_name")
	contactPhone := c.Query("contact_phone")

	items, err := services.FindLuggageByUserInfo(guestName, contactPhone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "query luggage success",
		"items":   items,
	})
}

// QueryLuggageByCode 按取件码查询寄存记录
// GET /storage/by-code?code=XXXX
func QueryLuggageByCode(c *gin.Context) {
	code := c.Query("code")
	item, err := services.FindLuggageByCode(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "query luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "query luggage success",
		"item":    item,
	})
}

// RetrieveLuggage 取件接口（通过取件码）
func RetrieveLuggage(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`         // 取件码
		RetrievedBy int64  `json:"retrieved_by" binding:"required"` // 操作员ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	item, err := services.RetrieveLuggage(req.Code, req.RetrievedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "retrieve luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "retrieve luggage success",
		"luggage_id": item.ID,
		"status":     item.Status,
	})
}

// ListLuggageByUser 获取用户寄存单列表
// GET /storage/list?user_id=1&status=stored
func ListLuggageByUser(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id is required",
		})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user_id",
		})
		return
	}

	status := c.Query("status")
	items, err := services.ListLuggageByUser(userID, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list luggage success",
		"items":   items,
	})
}

// ListLuggageByGuest 按客人姓名/手机号查询寄存单列表
// GET /storage/list/by-guest?guest_name=...&contact_phone=...&status=stored
func ListLuggageByGuest(c *gin.Context) {
	guestName := c.Query("guest_name")
	contactPhone := c.Query("contact_phone")
	status := c.Query("status")

	items, err := services.ListLuggageByGuest(guestName, contactPhone, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list luggage success",
		"items":   items,
	})
}

// GetLuggageDetail 获取寄存单详情
// GET /storage/detail?id=1
func GetLuggageDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	item, err := services.GetLuggageDetail(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get luggage detail failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get luggage detail success",
		"item":    item,
	})
}

// GetLuggageDetailByCode 按取件码查询寄存单详情
// GET /storage/detail/by-code?code=XXXX
func GetLuggageDetailByCode(c *gin.Context) {
	code := c.Query("code")
	item, err := services.GetLuggageDetailByCode(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get luggage detail failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get luggage detail success",
		"item":    item,
	})
}

// ListLuggageDetailByPhone 按手机号查询寄存单详情列表
// GET /storage/detail/by-phone?contact_phone=...
func ListLuggageDetailByPhone(c *gin.Context) {
	phone := c.Query("contact_phone")
	items, err := services.ListLuggageDetailByPhone(phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get luggage detail failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get luggage detail success",
		"items":   items,
	})
}
