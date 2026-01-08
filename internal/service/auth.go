package service

import (
	"errors"
	"time"

	"github.com/FruitsAI/Orange/internal/dto"
	"github.com/FruitsAI/Orange/internal/models"
	"github.com/FruitsAI/Orange/internal/pkg/jwt"
	"github.com/FruitsAI/Orange/internal/pkg/password"
	"github.com/FruitsAI/Orange/internal/repository"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// Login 用户登录
func (s *AuthService) Login(username, pwd string) (*dto.LoginResult, error) {
	// 查找用户
	user, err := s.userRepo.FindByCredential(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !password.CheckPassword(pwd, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账户已被禁用")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	// 更新最后登录时间
	now := time.Now()
	s.userRepo.UpdateFields(user.ID, map[string]interface{}{
		"last_login_time": now,
	})

	return &dto.LoginResult{Token: token, User: user}, nil
}

// Register 用户注册
func (s *AuthService) Register(input dto.RegisterRequest) error {
	// 检查用户名是否已存在
	if s.userRepo.ExistsByUsername(input.Username) {
		return errors.New("用户名已被注册")
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if input.Email != "" && s.userRepo.ExistsByEmail(input.Email) {
		return errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user := &models.User{
		Username: input.Username,
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: hashedPassword,
		Role:     "user",
		Status:   1,
	}

	return s.userRepo.Create(user)
}

// GetCurrentUser 获取当前用户
func (s *AuthService) GetCurrentUser(userID int64) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfile 更新个人信息
func (s *AuthService) UpdateProfile(userID int64, name, email, phone, department, position string) (*models.User, error) {
	updates := map[string]interface{}{}

	if name != "" {
		updates["name"] = name
	}
	if email != "" {
		updates["email"] = email
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if department != "" {
		updates["department"] = department
	}
	if position != "" {
		updates["position"] = position
	}

	if len(updates) > 0 {
		if err := s.userRepo.UpdateFields(userID, updates); err != nil {
			return nil, err
		}
	}

	return s.userRepo.FindByID(userID)
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !password.CheckPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := password.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	return s.userRepo.UpdateFields(userID, map[string]interface{}{
		"password": hashedPassword,
	})
}
