package services

import (
	"fmt"
	"os"
	"time"

	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/repository"
	"github.com/qujing226/pdf-enhancer/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.IUserRepository
}

// NewUserService 创建新的用户服务
func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

// GetUserByEmail 根据邮箱获取用户（包括密码哈希和盐值，用于认证）
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

// VerifyCredentials 验证用户凭据
func (s *UserService) VerifyCredentials(email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err // 用户不存在或查询错误
	}

	// 使用 utils.VerifyPassword 进行密码验证
	match, err := utils.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		// 哈希格式错误等内部问题
		return nil, fmt.Errorf("密码验证过程中发生错误: %w", err)
	}
	if !match {
		return nil, fmt.Errorf("密码错误")
	}

	// 移除敏感信息再返回
	user.PasswordHash = ""
	user.Salt = ""
	return user, nil
}

// GenerateToken 生成JWT令牌
func (s *UserService) GenerateToken(user *models.User) (string, error) {
	privateKeyPEM := os.Getenv("JWT_PRIVATE_KEY")
	if privateKeyPEM == "" {
		return "", fmt.Errorf("未配置JWT私钥")
	}
	return utils.GenerateJWT(user.ID, user.Email, privateKeyPEM)
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(name, email, password string) (*models.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 使用 utils.GeneratePasswordHash 生成密码哈希
	hashedPassword, err := utils.GeneratePasswordHash(password, nil)
	if err != nil {
		return nil, fmt.Errorf("生成密码哈希失败: %w", err)
	}

	// 创建新用户
	user := &models.User{
		ID:           utils.GenerateSnowflakeID(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 调用仓储层创建用户
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = "" // 不返回哈希
	user.Salt = ""         // 不返回盐
	return user, nil
}

// HashPassword 使用 bcrypt 哈希密码 (备用，推荐Argon2)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 检查密码和哈希是否匹配 (备用，推荐Argon2)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
