package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID           string    `json:"user_id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash,omitempty" db:"password_hash"`
	Salt         string    `json:"salt,omitempty" db:"salt"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Report 报告模型
type Report struct {
	ID        string    `json:"report_id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Summary   string    `json:"summary" db:"summary"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	PDFPath   string    `json:"pdf_path" db:"pdf_path"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ReportListItem 报告列表项
type ReportListItem struct {
	ReportID   string    `json:"report_id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	HasSummary bool      `json:"has_summary"`
}

// ReportListResponse 报告列表响应
type ReportListResponse struct {
	Reports []ReportListItem `json:"reports"`
}

// ReportDetailResponse 报告详情响应
type ReportDetailResponse struct {
	Report
	PDFURL string `json:"pdf_url"`
}

// SummaryRequest 摘要请求
type SummaryRequest struct {
	ReportID string `json:"report_id" binding:"required"`
}

// SummaryResponse 摘要响应
type SummaryResponse struct {
	Summary string `json:"summary"`
}

// APIResponse API通用响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewAPIResponse 创建新的API响应
func NewAPIResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
