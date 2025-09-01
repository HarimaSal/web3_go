package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

// errorHandler 全局错误处理中间件
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()
			var status int
			var msg string

			switch lastErr.Err {
			case gorm.ErrRecordNotFound:
				status = http.StatusNotFound
				msg = "资源不存在"
			case jwt.ErrSignatureInvalid, jwt.ErrTokenExpired:
				status = http.StatusUnauthorized
				msg = "无效或过期的认证令牌"
			default:
				status = http.StatusInternalServerError
				msg = "服务器内部错误"
			}

			logger := log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lshortfile)
			logger.Printf("错误: %v, 时间: %v", lastErr.Err, time.Now().Format(time.RFC3339))
			logger.Printf("响应状态：%d，消息：%s", status, msg)
		}
	}
}
