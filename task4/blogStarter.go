package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1、连接数据库
	db, err := initGormDB()
	if err != nil {
		fmt.Printf("init Gorm DB failed, err:%v\n", err)
		return
	}
	// 自动迁移模型
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		fmt.Printf("auto migrate failed, err:%v\n", err)
		return
	}

	// 2、创建路由
	router := gin.Default()
	//  注册
	router.POST("/register", func(c *gin.Context) {
		register(c, db)
	})
	// 登录
	router.POST("/login", func(c *gin.Context) {
		login(c, db)
	})

	/* 文章管理功能 */
	artGrp := router.Group("/article")
	artGrp.Use(AuthMiddleware()) // jwt验证中间件
	{
		// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
		artGrp.POST("/create", func(c *gin.Context) {

		})
		artGrp.GET("/listAll", func(c *gin.Context) {

		})
		// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
		// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
		// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
	}
	/* 评论功能 */

	/* 错误处理与日志记录 */

}
