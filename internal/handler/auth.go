package handler

import (
	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/pkg/response"
	"github.com/FruitsAI/Orange/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "用户名和密码不能为空")
		return
	}

	result, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Error(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// Register 用户注册
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	err := h.authService.Register(req)
	if err != nil {
		response.ParamError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "注册成功", nil)
}

// Logout 退出登录
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessWithMessage(c, "退出成功", nil)
}

// GetCurrentUser 获取当前用户信息
// GET /api/v1/users/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c)
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新个人信息
// PUT /api/v1/users/me
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID, req.Name, req.Email, req.Phone, req.Department, req.Position)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// ChangePassword 修改密码
// PUT /api/v1/users/me/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Unauthorized(c)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}
