package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

/* --- 评论板块 --- */

// 评论功能区
func commentBlock(router *gin.Engine, db *gorm.DB) {
	cmtGrp := router.Group("/comment")
	{
		// 实现评论的读取功能，支持获取某篇文章的所有评论列表。
		cmtGrp.GET("/:pid", func(c *gin.Context) {
			getAllCommentByPid(c, db)
		})
	}
	// 需认证的路由
	cmtAuthGrp := router.Group("/comment")
	cmtAuthGrp.Use(AuthMiddleware())
	{
		//实现评论的创建功能，已认证的用户可以对文章发表评论。
		cmtAuthGrp.POST("/create", func(c *gin.Context) {
			createComment(c, db)
		})
	}
}

// 创建评论。
func createComment(c *gin.Context, db *gorm.DB) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.Error(fmt.Errorf("评论输入无效"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论输入无效"})
		return
	}
	if comment.Content == "" {
		c.Error(fmt.Errorf("评论不能为空"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论不能为空"})
		return
	}
	// 从 JWT 中获取当前登录用户 ID
	userID := c.GetUint("userID")
	comment.UserID = userID
	if err := db.Create(&comment).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment": comment})
}

// 获取某篇文章的所有评论
func getAllCommentByPid(c *gin.Context, db *gorm.DB) {
	var comments []Comment
	if err := db.Where("post_id = ?", c.Param("pid")).Find(&comments).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comments retrieved successfully", "comments": comments})
}
