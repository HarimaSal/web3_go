package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	// 初始化日志
	logger, err := initLogger("blog.log")
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	// 1、连接数据库
	db, err := initGormDB(logger)
	if err != nil {
		logger.Fatalf("数据库连接失败: %v", err)
		return
	}
	// 自动迁移模型
	//err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	//if err != nil {
	//	logger.Fatalf("自动迁移失败: %v", err)
	//	return
	//}

	// 2、创建路由
	router := gin.Default()
	router.Use(gin.Logger())   // 使用 Gin 的日志中间件
	router.Use(errorHandler()) // 添加全局错误处理中间件
	//  注册
	router.POST("/register", func(c *gin.Context) {
		register(c, db)
	})
	// 登录
	router.POST("/login", func(c *gin.Context) {
		login(c, db)
	})

	/* 文章管理功能 */
	postBlock(router, db)

	/* 评论功能 */
	commentBlock(router, db)
	/* 错误处理与日志记录 */
	err = router.Run(":8090")
	if err != nil {
		logger.Fatalf("服务器启动失败: %v", err)
		return
	}
}

// initLogger 初始化并返回日志记录器
func initLogger(logFilePath string) (*log.Logger, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 使用 defer 确保文件关闭
	defer func(logFile *os.File) {
		if err := logFile.Close(); err != nil {
			log.Printf("关闭日志文件失败: %v", err)
		}
	}(logFile)

	logger := log.New(logFile, "BLOG: ", log.LstdFlags|log.Lshortfile)
	return logger, nil
}
