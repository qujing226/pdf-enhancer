package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config 数据库配置
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConfig 从环境变量创建数据库配置
func NewConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "13306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		DBName:   getEnv("DB_NAME", "reports_db"),
	}
}

// Connect 连接到数据库
func Connect(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	// 尝试连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法打开数据库连接: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 验证连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	log.Printf("成功连接到数据库 %s:%s/%s", config.Host, config.Port, config.DBName)
	return db, nil
}

// 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
