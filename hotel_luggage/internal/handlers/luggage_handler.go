package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"hotel_luggage/configs"
	"hotel_luggage/internal/repositories"
	"hotel_luggage/internal/services"
	"hotel_luggage/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

// CreateLuggageRequest 行李寄存请求结构体
type CreateLuggageRequest struct {
	GuestName    string                     `json:"guest_name" binding:"required"` // 客人用户名
	ContactPhone string                     `json:"contact_phone"`                 // 联系电话
	ContactEmail string                     `json:"contact_email"`                 // 联系邮箱
	Description  string                     `json:"description"`                   // 行李描述
	Quantity     int                        `json:"quantity"`                      // 行李数量
	SpecialNotes string                     `json:"special_notes"`                 // 特殊备注
	PhotoURL     string                     `json:"photo_url"`                     // 照片URL（可选）
	PhotoURLs    []string                   `json:"photo_urls"`                    // 多张照片URL（可选）
	StoreroomID  int64                      `json:"storeroom_id"`                  // 寄存室ID（单件模式必填）
	StaffName    string                     `json:"staff_name"`                    // 操作员用户名（可选，不传则用登录账号）
	QRCodeURL    string                     `json:"qr_code_url"`                   // 二维码URL（可选）
	Items        []CreateLuggageItemRequest `json:"items"`                         // 多件行李（可选）
}

// CreateLuggageItemRequest 单件行李（用于多件寄存）
type CreateLuggageItemRequest struct {
	StoreroomID  int64    `json:"storeroom_id" binding:"required"`
	Description  string   `json:"description"`
	Quantity     int      `json:"quantity"`
	SpecialNotes string   `json:"special_notes"`
	PhotoURL     string   `json:"photo_url"`
	PhotoURLs    []string `json:"photo_urls"`
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
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if req.StaffName == "" {
		req.StaffName = userNameStr
	}
	if req.StaffName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}
	if len(req.Items) > 0 {
		sharedCode := ""
		for i := 0; i < 5; i++ {
			codeCandidate, err := utils.GenerateCode(6)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "generate retrieval code failed",
					"error":   err.Error(),
				})
				return
			}
			exists, err := repositories.RetrievalCodeExists(codeCandidate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "generate retrieval code failed",
					"error":   err.Error(),
				})
				return
			}
			if !exists {
				sharedCode = codeCandidate
				break
			}
		}
		if sharedCode == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "generate retrieval code failed",
				"error":   "failed to generate unique retrieval code",
			})
			return
		}

		items := make([]gin.H, 0, len(req.Items))
		for _, it := range req.Items {
			created, err := services.CreateLuggage(services.CreateLuggageRequest{
				GuestName:     req.GuestName,
				ContactPhone:  req.ContactPhone,
				ContactEmail:  req.ContactEmail,
				Description:   it.Description,
				Quantity:      it.Quantity,
				SpecialNotes:  it.SpecialNotes,
				PhotoURL:      it.PhotoURL,
				PhotoURLs:     it.PhotoURLs,
				RetrievalCode: sharedCode,
				StoreroomID:   it.StoreroomID,
				StaffName:     req.StaffName,
				QRCodeURL:     req.QRCodeURL,
			})
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "create luggage failed",
					"error":   err.Error(),
				})
				return
			}
			items = append(items, gin.H{
				"luggage_id":   created.ID,
				"storeroom_id": created.StoreroomID,
				"photo_url":    created.PhotoURL,
				"photo_urls":   created.PhotoURLs,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message":        "create luggage success",
			"retrieval_code": sharedCode,
			"items":          items,
		})
		return
	}
	if req.StoreroomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   "storeroom_id is required",
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
		PhotoURL:     req.PhotoURL,
		PhotoURLs:    req.PhotoURLs,
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
		"photo_url":      item.PhotoURL,
		"photo_urls":     item.PhotoURLs,
	})
}

// QueryLuggageByUserInfo 按用户信息查询寄存记录
// GET /api/luggage/search?guest_name=...&contact_phone=...
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
// GET /api/luggage/by_code?code=XXXX
func QueryLuggageByCode(c *gin.Context) {
	code := c.Query("code")
	items, err := services.FindLuggageByCode(code)
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

	items, err := services.RetrieveLuggage(req.Code, req.RetrievedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "retrieve luggage failed",
			"error":   err.Error(),
		})
		return
	}
	luggageIDs := make([]int64, 0, len(items))
	for _, item := range items {
		luggageIDs = append(luggageIDs, item.ID)
	}
	var singleID interface{} = nil
	if len(luggageIDs) == 1 {
		singleID = luggageIDs[0]
	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "retrieve luggage success",
		"retrieval_code":  req.Code,
		"retrieved_by":    req.RetrievedBy,
		"retrieved_count": len(items),
		"luggage_ids":     luggageIDs,
		"luggage_id":      singleID,
	})
}

// CheckoutLuggageByCode 通过取件码取件
// POST /api/luggage/:id/checkout
func CheckoutLuggageByCode(c *gin.Context) {
	code := c.Param("id")
	username, _ := c.Get("username")
	retrievedBy, _ := username.(string)
	if retrievedBy == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}

	items, err := services.RetrieveLuggage(code, retrievedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "checkout failed",
			"error":   err.Error(),
		})
		return
	}
	luggageIDs := make([]int64, 0, len(items))
	for _, item := range items {
		luggageIDs = append(luggageIDs, item.ID)
	}
	var singleID interface{} = nil
	if len(luggageIDs) == 1 {
		singleID = luggageIDs[0]
	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "checkout success",
		"retrieval_code":  code,
		"retrieved_count": len(items),
		"luggage_ids":     luggageIDs,
		"luggage_id":      singleID,
	})
}

// GetCheckoutInfoByCode 获取当前酒店有行李在存的客人名单
// GET /api/luggage/:id/checkout
func GetCheckoutInfoByCode(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}

	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}

	items, err := services.ListGuestNamesByHotelAndStatus(*user.HotelID, "stored")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get checkout info failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "get checkout info success",
		"items":   items,
	})
}

// ListLuggageByUser 获取当前酒店所有已存放行李的客人姓名（去重）
// GET /api/luggage/list
func ListLuggageByUser(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}

	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}

	items, err := services.ListGuestNamesByHotelAndStatus(*user.HotelID, "stored")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list guest names failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list guest names success",
		"items":   items,
	})
}

// ListLuggageByGuest 按客人姓名/手机号查询寄存单列表
// GET /api/luggage/list/by_guest?guest_name=...&contact_phone=...&status=stored
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

// ListStoredLuggageByGuestName 查询某客人正在寄存的行李
// GET /api/luggage/list/by_guest_name?guest_name=...
func ListStoredLuggageByGuestName(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}
	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}

	guestName := c.Query("guest_name")
	items, err := services.ListStoredLuggageByGuestName(*user.HotelID, guestName)
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
// GET /api/luggage/:id
func GetLuggageDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("id")
	}
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
// GET /api/luggage/detail/by_code?code=XXXX
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
// GET /api/luggage/detail/by_phone?contact_phone=...
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
// GET /api/luggage/pickup_codes?username=xxx&status=stored
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
// GET /api/luggage/pickup_codes/by_phone?contact_phone=...&status=stored
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
// PUT /api/luggage/:id
func UpdateLuggageInfo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid luggage id",
		})
		return
	}

	var req struct {
		GuestName    *string   `json:"guest_name"`
		ContactPhone *string   `json:"contact_phone"`
		ContactEmail *string   `json:"contact_email"`
		Description  *string   `json:"description"`
		Quantity     *int      `json:"quantity"`
		SpecialNotes *string   `json:"special_notes"`
		PhotoURL     *string   `json:"photo_url"`
		PhotoURLs    *[]string `json:"photo_urls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
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
		PhotoURL:     req.PhotoURL,
		PhotoURLs:    req.PhotoURLs,
		UpdatedBy:    userNameStr,
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
// PUT /api/luggage/:id/code
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
// POST /api/luggage/bind
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

// Upload 上传图片接口（multipart/form-data）
// POST /api/upload
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "missing file",
		})
		return
	}
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "file too large",
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "invalid file extension",
		})
		return
	}
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "upload failed",
			"error":   "unsupported file type",
		})
		return
	}

	now := time.Now()
	// 生成随机文件名
	nameBytes := make([]byte, 16)
	if _, err := rand.Read(nameBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "upload failed",
			"error":   err.Error(),
		})
		return
	}
	fileName := hex.EncodeToString(nameBytes) + ext
	objectName := fmt.Sprintf("uploads/%s/%s/%s", now.Format("2006"), now.Format("01"), fileName)

	// 优先使用MinIO上传
	if repositories.MinIOClient != nil {
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "upload failed",
				"error":   "cannot open file",
			})
			return
		}
		defer fileReader.Close()

		contentType := file.Header.Get("Content-Type")
		if contentType == "" {
			contentType = mime.TypeByExtension(ext)
		}
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// 设置30秒超时（考虑网络延迟和文件大小）
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		
		_, err = repositories.MinIOClient.PutObject(ctx, repositories.MinIOBucketName, objectName, fileReader, file.Size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			// MinIO上传失败，降级到本地存储
			log.Printf("⚠️  MinIO上传失败(超时或网络错误)，降级到本地存储: %v", err)
		} else {
			// MinIO上传成功，返回MinIO URL
			minioConfig := configs.LoadMinIOConfig()
			scheme := "http"
			if minioConfig.UseSSL {
				scheme = "https"
			}
			fullURL := fmt.Sprintf("%s://%s/%s/%s", scheme, minioConfig.Endpoint, repositories.MinIOBucketName, objectName)

			c.JSON(http.StatusOK, gin.H{
				"message":       "upload success (MinIO)",
				"url":           fullURL,
				"object_name":   objectName,
				"content_type":  contentType,
				"size":          file.Size,
				"file_name":     fileName,
				"max_size_byte": maxSize,
				"storage":       "minio",
			})
			return
		}
	}

	// 降级方案：本地文件存储
	subDir := filepath.Join("uploads", now.Format("2006"), now.Format("01"))
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "upload failed",
			"error":   err.Error(),
		})
		return
	}

	savePath := filepath.Join(subDir, fileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "upload failed",
			"error":   err.Error(),
		})
		return
	}

	relativeURL := fmt.Sprintf("/uploads/%s/%s/%s", now.Format("2006"), now.Format("01"), fileName)
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fullURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, relativeURL)
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "upload success (Local)",
		"url":           fullURL,
		"relative_url":  relativeURL,
		"content_type":  contentType,
		"size":          file.Size,
		"file_name":     fileName,
		"max_size_byte": maxSize,
		"storage":       "local",
	})
}

// ListHistoryByGuest 查询取件历史（按客人姓名/手机号）
// GET /api/luggage/history?guest_name=...&contact_phone=...
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

// ListLuggageByStoreroom 获取寄存室下的行李订单列表
// GET /api/luggage/storerooms/:id/orders?status=stored
func ListLuggageByStoreroom(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid storeroom id",
		})
		return
	}
	status := c.Query("status")
	items, err := services.ListLuggageByStoreroom(id, status)
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

// ListStoredLogs 获取所有寄存记录（当前登录用户的酒店）
// GET /api/luggage/logs/stored
func ListStoredLogs(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}
	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}
	items, err := services.ListLuggageByHotelAndStatus(*user.HotelID, "stored")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   items,
	})
}

// ListUpdatedLogs 获取所有寄存信息修改记录（当前登录用户的酒店）
// GET /api/luggage/logs/updated
func ListUpdatedLogs(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}
	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}
	items, err := services.ListLuggageUpdatesByHotel(*user.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   items,
	})
}

// ListRetrievedLogs 获取所有取出记录（当前登录用户的酒店）
// GET /api/luggage/logs/retrieved
func ListRetrievedLogs(c *gin.Context) {
	username, _ := c.Get("username")
	userNameStr, _ := username.(string)
	if userNameStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "missing user info",
		})
		return
	}
	user, err := repositories.GetUserByUsername(userNameStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "get user failed",
			"error":   err.Error(),
		})
		return
	}
	if user.HotelID == nil || *user.HotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is missing",
		})
		return
	}
	items, err := services.ListHistoryByHotel(*user.HotelID, "", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list logs failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list logs success",
		"items":   items,
	})
}
