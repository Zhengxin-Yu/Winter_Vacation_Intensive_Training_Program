package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateStoreroomRequest 创建寄存室请求结构体
type CreateStoreroomRequest struct {
	Name     string `json:"name" binding:"required"` // 寄存室名称
	Location string `json:"location"`                // 位置描述
	Capacity int    `json:"capacity"`                // 容量
	IsActive bool   `json:"is_active"`               // 是否启用
}

// UpdateStoreroomStatusRequest 更新寄存室状态请求结构体
type UpdateStoreroomStatusRequest struct {
	IsActive bool `json:"is_active"` // 是否启用
}

// ListStorerooms 获取寄存室列表
func ListStorerooms(c *gin.Context) {
	rooms, err := services.ListStorerooms()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list storerooms failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list storerooms success",
		"items":   rooms,
	})
}

// CreateStoreroom 创建寄存室
func CreateStoreroom(c *gin.Context) {
	var req CreateStoreroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	room, err := services.CreateStoreroom(services.CreateStoreroomRequest{
		Name:     req.Name,
		Location: req.Location,
		Capacity: req.Capacity,
		IsActive: req.IsActive,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create storeroom failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create storeroom success",
		"item":    room,
	})
}

// DeleteStoreroom 删除寄存室（有行李不能删）
func DeleteStoreroom(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid storeroom id",
		})
		return
	}

	if err := services.DeleteStoreroom(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete storeroom failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete storeroom success",
	})
}

// UpdateStoreroomStatus 更新寄存室启用状态
func UpdateStoreroomStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid storeroom id",
		})
		return
	}

	var req UpdateStoreroomStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := services.UpdateStoreroomStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update storeroom status failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update storeroom status success",
	})
}
