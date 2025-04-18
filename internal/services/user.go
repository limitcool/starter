package services

import (
	"fmt"
	"time"

	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/repository"
)

// UserService 普通用户服务
type UserService struct {
	userRepo *repository.GormUserRepository
}

// NewUserService 创建普通用户服务
func NewUserService(userRepo *repository.GormUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// VerifyPassword 验证用户密码
func (s *UserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// Register 用户注册
func (s *UserService) Register(req v1.UserRegisterRequest, registerIP string) (*model.User, error) {
	isExist, err := s.userRepo.IsExist(req.Username)
	if err != nil {
		return nil, errorx.ErrDatabaseQueryError.WithError(err)
	}
	if isExist {
		return nil, errorx.ErrUserAlreadyExists
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errorx.ErrInternal.WithError(err)
	}

	// 创建用户
	user := &model.User{
		Username:   req.Username,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Mobile:     req.Mobile,
		Enabled:    true,
		Gender:     req.Gender,
		Birthday:   &time.Time{},
		Address:    req.Address,
		RegisterIP: registerIP,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errorx.ErrDatabaseInsertError.WithError(err)
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, errorx.ErrUserDisabled
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, errorx.ErrUserPasswordError.WithMsg("密码错误")
	}

	// 更新最后登录时间和IP
	fields := map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	}
	if err := s.userRepo.UpdateFields(user.ID, fields); err != nil {
		return nil, err
	}

	// 获取配置
	cfg := core.Instance().Config()

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeUser,             // 普通用户
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeUser,              // 普通用户
		TokenType: enum.TokenTypeRefresh.String(), // 刷新令牌
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := jwtpkg.GenerateTokenWithCustomClaims(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    cfg.JwtAuth.AccessExpire,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, data map[string]any) error {
	// 不允许更新的字段
	delete(data, "id")
	delete(data, "username")
	delete(data, "password")
	delete(data, "created_at")
	delete(data, "deleted_at")

	// 更新用户信息
	fields := make(map[string]interface{}, len(data))
	for k, v := range data {
		fields[k] = v
	}
	return s.userRepo.UpdateFields(int64(id), fields)
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id int64, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(oldPassword, user.Password) {
		return errorx.ErrUserPasswordError.WithMsg("原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	fields := map[string]interface{}{
		"password": hashedPassword,
	}
	return s.userRepo.UpdateFields(id, fields)
}

// 该函数已移动到 controller/user.go 中的 UserInfo 方法
