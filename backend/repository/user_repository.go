package repository

import (
	"database/sql"
	"fmt"

	"github.com/qujing226/pdf-enhancer/backend/models"
)

// IUserRepository  用户仓储接口
type IUserRepository interface {
	GetByID(userID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

// UserRepository 用户仓储实现
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(userID string) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`
	user := &models.User{}
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash, salt, created_at, updated_at FROM users WHERE email = ?`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Salt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return user, nil
}

// Create 创建新用户
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, name, email, password_hash,salt, created_at, updated_at) 
          VALUES (?, ?, ?,?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.PasswordHash, user.Salt, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}
