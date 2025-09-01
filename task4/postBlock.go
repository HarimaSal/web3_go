package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

/* --- 文章板块 --- */
// 文章管理功能区
func postBlock(router *gin.Engine, db *gorm.DB) {
	artGrp := router.Group("/article")
	{
		// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
		artGrp.GET("/listAll", func(c *gin.Context) {
			listAllPosts(c, db)
		})
		artGrp.GET("/:id", func(c *gin.Context) {
			getPost(c, db)
		})
	}
	// 需要权限认证的路由
	artAuthGrp := router.Group("/article")
	artAuthGrp.Use(AuthMiddleware()) // jwt验证中间件
	{
		// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
		artAuthGrp.POST("/create", func(c *gin.Context) {
			// 创建文章
			createPost(c, db)
		})
		// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
		artAuthGrp.PUT("/:id", func(c *gin.Context) {
			updatePost(c, db)
		})
		// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
		artAuthGrp.DELETE("/:id", func(c *gin.Context) {
			deletePost(c, db)
		})
	}
}

// 删除文章
func deletePost(c *gin.Context, db *gorm.DB) {
	// 从 JWT 中获取当前登录用户 ID
	userID := c.GetUint("userID")
	var post Post
	if err := db.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}
	// 权限检查：仅允许文章作者更新
	if post.UserID != userID {
		c.Error(fmt.Errorf("无权限更新该文章"))
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新该文章"})
		return
	}
	if err := db.Delete(&post).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// 更新文章
func updatePost(c *gin.Context, db *gorm.DB) {
	// 从 JWT 中获取当前登录用户 ID
	userID := c.GetUint("userID")
	// 获取文章 ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}
	// 检查文章是否存在
	var post Post
	if err := db.First(&post, id).Error; err != nil {
		c.Error(err) // 记录错误
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	// 权限检查：仅允许文章作者更新
	if post.UserID != userID {
		c.Error(fmt.Errorf("无权限更新该文章"))
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限更新该文章"})
		return
	}
	var uptPost Post
	// 使用 Post 结构体绑定请求数据
	if err := c.ShouldBindJSON(&uptPost); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 只更新允许的字段（title 和 content），避免修改 UserID 或其他字段
	if uptPost.Title != "" {
		post.Title = uptPost.Title
	}
	if uptPost.Content != "" {
		post.Content = uptPost.Content
	}
	// 保存更新
	if err := db.Save(&post).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章更新失败"})
		return
	}
}

// 获取文章
func getPost(c *gin.Context, db *gorm.DB) {
	var post Post
	if err := db.Preload("Comments").First(&post, 1).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve post"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// 获取所有文章
func listAllPosts(c *gin.Context, db *gorm.DB) {
	var posts []Post
	if err := db.Find(&posts).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// 创建文章
func createPost(c *gin.Context, db *gorm.DB) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 验证标题和内容
	if post.Title == "" || post.Content == "" {
		c.Error(fmt.Errorf("标题或内容为空"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "标题或内容为空"})
		return
	}
	// 创建文章入库
	if err := db.Create(&post).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文章创建失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": post})
}
