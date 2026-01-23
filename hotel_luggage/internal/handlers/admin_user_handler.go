package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// ListUsersByHotel 获取指定酒店的用户列表
// GET /api/luggage/users?hotel_id=1
func ListUsersByHotel(c *gin.Context) {
	hotelIDStr := c.Query("hotel_id")
	if hotelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hotel_id is required",
		})
		return
	}
	hotelID, err := strconv.ParseInt(hotelIDStr, 10, 64)
	if err != nil || hotelID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid hotel_id",
		})
		return
	}

	items, err := services.ListUsersByHotel(hotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list users failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list users success",
		"items":   items,
	})
}
