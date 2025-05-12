package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/qujing226/pdf-enhancer/backend/models"
	"github.com/qujing226/pdf-enhancer/backend/services"
)

// AuthHandler 处理认证相关的请求
type AuthHandler struct {
	authService *services.UserService
}

// NewAuthHandler 创建新的认证处理器
func NewAuthHandler(authService *services.UserService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 处理用户注册请求
func (h *AuthHandler) Register(c *gin.Context) {
	var registerReq models.RegisterRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(400, models.NewAPIResponse(400, "无效的请求参数", nil))
		return
	}

	// 创建新用户
	user, err := h.authService.CreateUser(registerReq.Name, registerReq.Email, registerReq.Password)
	if err != nil {
		c.JSON(400, models.NewAPIResponse(400, err.Error(), nil))
		return
	}

	// 生成JWT令牌
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(500, models.NewAPIResponse(500, "生成令牌失败", err.Error()))
		return
	}

	// 返回注册成功响应
	c.JSON(201, models.NewAPIResponse(201, "注册成功", models.RegisterResponse{
		Token: token,
		User:  *user,
	}))
}

// Login 处理用户登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq models.LoginRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, models.NewAPIResponse(400, "无效的请求参数", nil))
		return
	}

	// 验证用户凭据
	user, err := h.authService.VerifyCredentials(loginReq.Email, loginReq.Password)
	if err != nil {
		c.JSON(401, models.NewAPIResponse(401, "认证失败", err.Error()))
		return
	}

	// 生成JWT令牌
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(500, models.NewAPIResponse(500, "生成令牌失败", err.Error()))
		return
	}

	// 返回登录成功响应
	c.JSON(200, models.NewAPIResponse(200, "登录成功", models.LoginResponse{
		Token: token,
		User:  *user,
	}))
}
