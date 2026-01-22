package handlers

import (
	"net/http"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// MigrateLuggageRequest 行李迁移请求结构体
type MigrateLuggageRequest struct {
	LuggageID     int64  `json:"luggage_id" binding:"required"`      // 行李ID
	ToStoreroomID int64  `json:"to_storeroom_id" binding:"required"` // 目标寄存室ID
	MigratedBy    string `json:"migrated_by" binding:"required"`     // 操作员用户名
}

// MigrateLuggage 行李迁移接口
func MigrateLuggage(c *gin.Context) {
	var req MigrateLuggageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	log, err := services.MigrateLuggage(services.MigrateLuggageRequest{
		LuggageID:     req.LuggageID,
		ToStoreroomID: req.ToStoreroomID,
		MigratedBy:    req.MigratedBy,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "migrate luggage failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "migrate luggage success",
		"log":     log,
	})
}
