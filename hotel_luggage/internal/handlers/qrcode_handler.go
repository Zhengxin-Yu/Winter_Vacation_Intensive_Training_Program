package handlers

import (
	"net/http"

	"hotel_luggage/utils"

	"github.com/gin-gonic/gin"
)

// GetQRCode 返回二维码图片
// GET /qr/:code
func GetQRCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "code is required",
		})
		return
	}

	png, err := utils.GenerateQRCodePNG(code, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "generate qr failed",
			"error":   err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}
