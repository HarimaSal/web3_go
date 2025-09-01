package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 HTTP 请求头 "Authorization" 中获取 JWT Token
		// 常见的格式是 "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(fmt.Errorf("Authorization header is required"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort() // 中止后续处理
			return
		}

		// 2. 提取 Token 字符串（去除 "Bearer " 前缀）
		// 约定俗成的格式是 "Bearer <token>"，所以需要去除前7个字符
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.Error(fmt.Errorf("Invalid token format"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:] // 提取实际的 Token 部分

		// 3. 解析并验证 JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法是否为 HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 返回密钥
			return []byte(secretKey), nil
		})

		// 4. 处理解析或验证过程中的错误
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// 5. 检查 Token 是否有效
		if !token.Valid {
			c.Error(fmt.Errorf("Invalid token"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 6. 从已验证的 Token 中提取 Claims（声明信息）
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// 将用户信息（如用户ID）存储到 Gin 的上下文中，以便后续处理函数使用
			// 提取并转换 userID
			if id, ok := claims["id"].(float64); ok {
				c.Set("userID", uint(id)) // 转换为 uint
			} else {
				c.Error(fmt.Errorf("invalid user ID in token"))
				c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户 ID"})
				c.Abort()
				return
			}
			c.Set("username", claims["username"])
			// 你可以根据需要存储其他声明信息
		} else {
			c.Error(fmt.Errorf("Failed to parse token claims"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token claims"})
			c.Abort()
			return
		}
		// 7. 验证通过，继续执行后续的请求处理链（例如最终的业务处理函数）
		c.Next()
	}
}

// 注册
func register(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err) // 记录错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser User
	if err := db.Where("username = ? ", user.Username).First(&existingUser).Error; err == nil {
		c.Error(fmt.Errorf("用户名已存在"))
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	// 用户存入库
	if err := db.Create(&user).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// 签名密钥
var secretKey = "sdjhakdhajdklsl;0653632"

// 登录
func login(c *gin.Context, db *gorm.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询用户
	var existingUser User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名错误"})
		return
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       existingUser.ID,
		"username": existingUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT 生成失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tokenString": tokenString})
}
