package handlers

import (
	"net/http"

	"hotel_luggage/internal/services"
	"hotel_luggage/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 明文密码
}

// CreateUserRequest 创建用户请求结构体
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 明文密码
	Role     string `json:"role"`                        // 角色，可选：staff/admin
	HotelID  *int64 `json:"hotel_id"`                    // 关联酒店ID（staff/admin 必填）
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	// 解析并校验 JSON 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	user, err := services.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "login failed",
			"error":   err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "token generate failed",
			"error":   err.Error(),
		})
		return
	}

	// 目前先返回基础信息，后续可接入 JWT
	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"hotel_id": user.HotelID,
		},
		"token": token,
	})
}

// CreateUser 创建用户接口（自动生成 bcrypt 密码哈希）
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	// 解析并校验 JSON 参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req.Username, req.Password, req.Role, req.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create user failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create user success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"hotel_id": user.HotelID,
		},
	})
}
