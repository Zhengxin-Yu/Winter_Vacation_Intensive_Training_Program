package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateHotelRequest 创建酒店请求
type CreateHotelRequest struct {
	Name     string `json:"name" binding:"required"` // 酒店名称
	Address  string `json:"address"`                 // 地址
	Phone    string `json:"phone"`                   // 电话
	IsActive bool   `json:"is_active"`               // 是否启用
}

// UpdateHotelRequest 更新酒店请求
type UpdateHotelRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

// ListHotels 获取酒店列表
func ListHotels(c *gin.Context) {
	items, err := services.ListHotels()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list hotels failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list hotels success",
		"items":   items,
	})
}

// CreateHotel 创建酒店
func CreateHotel(c *gin.Context) {
	var req CreateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	hotel, err := services.CreateHotel(req.Name, req.Address, req.Phone, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create hotel failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "create hotel success",
		"item":    hotel,
	})
}

// UpdateHotel 更新酒店
func UpdateHotel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid hotel id",
		})
		return
	}

	var req UpdateHotelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	if err := services.UpdateHotel(id, req.Name, req.Address, req.Phone, req.IsActive); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "update hotel failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update hotel success",
	})
}

// DeleteHotel 删除酒店
func DeleteHotel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid hotel id",
		})
		return
	}

	if err := services.DeleteHotel(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete hotel failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete hotel success",
	})
}
