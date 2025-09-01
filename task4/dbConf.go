package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func initGormDB(*log.Logger) (*gorm.DB, error) {
	// 数据库连接配置
	dsn := "root:root@tcp(localhost:3308)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库表创建成功！")

	return gormDb, err
}
