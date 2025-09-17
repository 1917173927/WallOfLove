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

var DB *gorm.DB

func Init() {
	host := config.Config.GetString("mysql.host")
	port := config.Config.GetString("mysql.port")
	user := config.Config.GetString("mysql.user")
	password := config.Config.GetString("mysql.password")
	DBname := config.Config.GetString("mysql.DBname")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, DBname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	tx, err := sqlDB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	err = db.WithContext(context.Background()).
		Session(&gorm.Session{NewDB: true}).
		AutoMigrate(&models.User{}, &models.Post{}, &models.Image{}, &models.Blacklist{}, &models.Review{}, &models.Review2{})
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
