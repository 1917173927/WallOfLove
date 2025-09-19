// Package database 提供数据库连接和初始化的功能
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/1917173927/WallOfLove/app/models"
	"github.com/1917173927/WallOfLove/conf/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 是全局数据库连接实例
var DB *gorm.DB

// Init 初始化数据库连接并执行自动迁移
// 从配置文件中读取数据库连接信息，并创建全局数据库实例
func Init() {
	// 从配置文件中读取数据库连接信息
	host := config.Config.GetString("mysql.host")
	port := config.Config.GetString("mysql.port")
	user := config.Config.GetString("mysql.user")
	password := config.Config.GetString("mysql.password")
	DBname := config.Config.GetString("mysql.DBname")

	// 构建数据库连接字符串 (DSN)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, DBname)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// 获取底层 SQL 数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	// 开始数据库事务
	tx, err := sqlDB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	// 执行自动迁移，确保数据库表结构与模型定义一致
	err = db.WithContext(context.Background()).
		Session(&gorm.Session{NewDB: true}).
		AutoMigrate(&models.User{}, &models.Post{}, &models.Image{}, &models.Blacklist{}, &models.Review{}, &models.Reply{})
	if err != nil {
		_ = tx.Rollback()
		log.Fatal("自动迁移失败:", err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal("事务提交失败:", err)
	}
	DB = db
	log.Println("数据库初始化完成")
}
