package dao_models

import (
	"time"
)

// ReportDAO 报告数据库模型
type ReportDAO struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"` // 使用longtext类型存储
	Summary   string    `db:"summary"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	PDFPath   string    `db:"pdf_path"`
}

// UserDAO 用户数据库模型
type UserDAO struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Salt         string    `db:"salt"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}