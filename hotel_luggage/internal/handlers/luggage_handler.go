package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateLuggageRequest 行李寄存请求结构体
type CreateLuggageRequest struct {
	GuestName    string `json:"guest_name" binding:"required"`   // 客人用户名
	ContactPhone string `json:"contact_phone"`                   // 联系电话
	ContactEmail string `json:"contact_email"`                   // 联系邮箱
	Description  string `json:"description"`                     // 行李描述
	Quantity     int    `json:"quantity"`                        // 行李数量
	SpecialNotes string `json:"special_notes"`                   // 特殊备注
	StoreroomID  int64  `json:"storeroom_id" binding:"required"` // 寄存室ID
	StaffName    string `json:"staff_name" binding:"required"`   // 操作员用户名
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
		StaffName:    req.StaffName,
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
		"qrcode_url":     item.QRCodeURL,
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

// QueryLuggageByPhone 按手机号查询寄存记录
// GET /api/luggage/by_phone?contact_phone=...
func QueryLuggageByPhone(c *gin.Context) {
	phone := c.Query("contact_phone")
	items, err := services.FindLuggageByUserInfo("", phone)
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

// RetrieveLuggage 取件接口（通过取件码）
func RetrieveLuggage(c *gin.Context) {
	var req struct {
		Code        string `json:"code" binding:"required"`         // 取件码
		RetrievedBy string `json:"retrieved_by" binding:"required"` // 操作员用户名
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
		"message":       "retrieve luggage success",
		"luggage_id":    item.ID,
		"status":        item.Status,
		"retrieved_by":  req.RetrievedBy,
	})
}

// CheckoutLuggageByCode 通过取件码取件
// POST /api/luggage/:id/checkout
func CheckoutLuggageByCode(c *gin.Context) {
	code := c.Param("id")
	var req struct {
		RetrievedBy string `json:"retrieved_by" binding:"required"` // 操作员用户名
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	item, err := services.RetrieveLuggage(code, req.RetrievedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "checkout failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "checkout success",
		"luggage_id": item.ID,
	})
}

// ListLuggageByUser 获取用户寄存单列表
// GET /storage/list?username=xxx&status=stored
func ListLuggageByUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username is required",
		})
		return
	}

	status := c.Query("status")
	items, err := services.ListLuggageByUser(username, status)
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

// ListPickupCodesByUser 查看取件码列表
// GET /pickup-codes?username=xxx&status=stored
func ListPickupCodesByUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username is required",
		})
		return
	}

	status := c.Query("status")
	items, err := services.ListPickupCodesByUser(username, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get pickup codes failed",
			"error":   err.Error(),
		})
		return
	}

	// 只返回取件码相关字段，避免暴露不必要数据
	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		result = append(result, gin.H{
			"luggage_id":     item.ID,
			"guest_name":     item.GuestName,
			"contact_phone":  item.ContactPhone,
			"retrieval_code": item.RetrievalCode,
			"status":         item.Status,
			"stored_at":      item.StoredAt,
			"retrieved_at":   item.RetrievedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get pickup codes success",
		"items":   result,
	})
}

// ListPickupCodesByPhone 按客人手机号查询取件码列表
// GET /pickup-codes/by-phone?contact_phone=...&status=stored
func ListPickupCodesByPhone(c *gin.Context) {
	phone := c.Query("contact_phone")
	status := c.Query("status")
	items, err := services.ListPickupCodesByPhone(phone, status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get pickup codes failed",
			"error":   err.Error(),
		})
		return
	}

	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		result = append(result, gin.H{
			"luggage_id":     item.ID,
			"guest_name":     item.GuestName,
			"contact_phone":  item.ContactPhone,
			"retrieval_code": item.RetrievalCode,
			"status":         item.Status,
			"stored_at":      item.StoredAt,
			"retrieved_at":   item.RetrievedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get pickup codes success",
		"items":   result,
	})
}

// UpdateLuggageInfo 修改寄存信息（不包含寄存室迁移）
// PUT /storage/:id
func UpdateLuggageInfo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid luggage id",
		})
		return
	}

	var req struct {
		GuestName    *string `json:"guest_name"`
		ContactPhone *string `json:"contact_phone"`
		ContactEmail *string `json:"contact_email"`
		Description  *string `json:"description"`
		Quantity     *int    `json:"quantity"`
		SpecialNotes *string `json:"special_notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := services.UpdateLuggageInfo(id, services.UpdateLuggageInfoRequest{
		GuestName:    req.GuestName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Description:  req.Description,
		Quantity:     req.Quantity,
		SpecialNotes: req.SpecialNotes,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update luggage success",
	})
}

// UpdateLuggageCode 修改取件码
// PUT /storage/:id/code
func UpdateLuggageCode(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid luggage id",
		})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := services.UpdateLuggageCode(id, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update retrieval code failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update retrieval code success",
	})
}

// BindLuggage 绑定行李到用户
// POST /storage/bind
func BindLuggage(c *gin.Context) {
	var req struct {
		LuggageID int64  `json:"luggage_id" binding:"required"` // 行李ID
		Username  string `json:"username" binding:"required"`   // 用户名
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := services.BindLuggageToUser(req.LuggageID, req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bind luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bind luggage success",
	})
}

// TransferLuggageByID 迁移行李
// POST /api/luggage/:id/transfer
func TransferLuggageByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid luggage id",
		})
		return
	}

	var req struct {
		ToStoreroomID int64  `json:"to_storeroom_id" binding:"required"`
		MigratedBy    string `json:"migrated_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	_, err = services.MigrateLuggage(services.MigrateLuggageRequest{
		LuggageID:     id,
		ToStoreroomID: req.ToStoreroomID,
		MigratedBy:    req.MigratedBy,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "transfer failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transfer success",
	})
}

// ListTransfersByLuggageID 查询迁移历史
// GET /api/luggage/:id/transfers
func ListTransfersByLuggageID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid luggage id",
		})
		return
	}

	items, err := services.ListMigrationsByLuggageID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list transfers failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list transfers success",
		"items":   items,
	})
}

// Upload 占位接口（上传功能）
// POST /api/upload
func Upload(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "upload not implemented",
	})
}

// ListHistoryByGuest 查询取件历史（按客人姓名/手机号）
// GET /storage/history/by-guest?guest_name=...&contact_phone=...
func ListHistoryByGuest(c *gin.Context) {
	guestName := c.Query("guest_name")
	contactPhone := c.Query("contact_phone")

	items, err := services.ListHistoryByGuest(guestName, contactPhone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get history failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "get history success",
		"items":   items,
	})
}
