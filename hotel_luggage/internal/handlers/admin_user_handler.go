package handlers

import (
	"net/http"
	"strconv"

	"hotel_luggage/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateEmployeeRequest 创建员工请求
type CreateEmployeeRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	HotelID  *int64 `json:"hotel_id" binding:"required"`
}

// CreateAdminRequest 创建管理员请求
type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	HotelID  *int64 `json:"hotel_id" binding:"required"`
}

// CreateEmployee 创建员工
func CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req.Username, req.Password, "staff", req.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create employee failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create employee success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"hotel_id": user.HotelID,
		},
	})
}

// CreateAdmin 创建管理员
func CreateAdmin(c *gin.Context) {
	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
			"error":   err.Error(),
		})
		return
	}

	user, err := services.CreateUser(req.Username, req.Password, "admin", req.HotelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "create admin failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create admin success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
			"hotel_id": user.HotelID,
		},
	})
}

// ListEmployees 员工列表
func ListEmployees(c *gin.Context) {
	items, err := services.ListUsersByRole("staff")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list employees failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list employees success",
		"items":   items,
	})
}

// ListAdmins 管理员列表
func ListAdmins(c *gin.Context) {
	items, err := services.ListUsersByRole("admin")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "list admins failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "list admins success",
		"items":   items,
	})
}

// DeleteEmployee 删除员工
func DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := services.DeleteUserByRole(id, "staff"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete employee failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete employee success",
	})
}

// DeleteAdmin 删除管理员
func DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	if err := services.DeleteUserByRole(id, "admin"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "delete admin failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete admin success",
	})
}

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
