package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/repositories"
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

// ListStorerooms 获取寄存室列表（当前登录用户的酒店）
func ListStorerooms(c *gin.Context) {
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

	rooms, err := services.ListStorerooms(*user.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list storerooms failed",
			"error":   err.Error(),
		})
		return
	}

	result := make([]gin.H, 0, len(rooms))
	for _, room := range rooms {
		storedCount, err := repositories.CountStoredByStoreroom(room.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "count storeroom luggage failed",
				"error":   err.Error(),
			})
			return
		}
		remaining := -1
		if room.Capacity > 0 {
			remaining = room.Capacity - int(storedCount)
			if remaining < 0 {
				remaining = 0
			}
		}
		result = append(result, gin.H{
			"id":                 room.ID,
			"hotel_id":           room.HotelID,
			"name":               room.Name,
			"location":           room.Location,
			"capacity":           room.Capacity,
			"is_active":          room.IsActive,
			"created_at":         room.CreatedAt,
			"stored_count":       storedCount,
			"remaining_capacity": remaining,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "list storerooms success",
		"items":   result,
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

	room, err := services.CreateStoreroom(services.CreateStoreroomRequest{
		HotelID:  *user.HotelID,
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
