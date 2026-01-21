package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home 返回系统功能入口（给前端提供基础导航信息）
// GET /home
func Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hotel luggage system api",
		"endpoints": []gin.H{
			{"name": "健康检查", "method": "GET", "path": "/ping"},
			{"name": "登录", "method": "POST", "path": "/login"},
			{"name": "创建用户", "method": "POST", "path": "/users"},
			{"name": "行李寄存", "method": "POST", "path": "/storage"},
			{"name": "取件", "method": "POST", "path": "/storage/retrieve"},
			{"name": "寄存查询（姓名/手机号）", "method": "GET", "path": "/storage/search"},
			{"name": "寄存查询（取件码）", "method": "GET", "path": "/storage/by-code"},
			{"name": "寄存单列表（用户）", "method": "GET", "path": "/storage/list"},
			{"name": "寄存单列表（客人）", "method": "GET", "path": "/storage/list/by-guest"},
			{"name": "寄存单详情（ID）", "method": "GET", "path": "/storage/detail"},
			{"name": "寄存单详情（取件码）", "method": "GET", "path": "/storage/detail/by-code"},
			{"name": "寄存单详情（手机号）", "method": "GET", "path": "/storage/detail/by-phone"},
			{"name": "取件码列表（用户）", "method": "GET", "path": "/pickup-codes"},
			{"name": "取件码列表（手机号）", "method": "GET", "path": "/pickup-codes/by-phone"},
			{"name": "寄存室列表", "method": "GET", "path": "/storerooms"},
			{"name": "创建寄存室", "method": "POST", "path": "/storerooms"},
			{"name": "删除寄存室", "method": "DELETE", "path": "/storerooms/:id"},
			{"name": "寄存室状态更新", "method": "PUT", "path": "/storerooms/:id/status"},
			{"name": "行李迁移", "method": "POST", "path": "/storerooms/migrate"},
			{"name": "二维码展示", "method": "GET", "path": "/qr/:code"},
		},
	})
}
