package utils

import "github.com/skip2/go-qrcode"

// GenerateQRCodePNG 生成二维码 PNG 字节
// content 为二维码内容，size 为图片尺寸（像素）
func GenerateQRCodePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	return qrcode.Encode(content, qrcode.Medium, size)
}
