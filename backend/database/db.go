package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Zhaobo-Wang/go-projects/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if username == "" {
		username = "root"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}
	if dbname == "" {
		dbname = "todo_app"
	}

	// 1) 先连接到不带数据库名的 MySQL（用于创建数据库）
	dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port)
	tmpDB, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL server:", err)
	}

	// 创建数据库（如果不存在）
	if err := tmpDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbname).Error; err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// 2) 再用带数据库名的 DSN 打开连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 3) 自动迁移（把错误打印并在出错时终止，方便发现问题）
	if err := DB.AutoMigrate(&models.User{}, &models.Todo{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	// 可选：检查表是否存在，便于调试
	has := DB.Migrator().HasTable(&models.Todo{})
	log.Printf("Has todos table? %v", has)
}
