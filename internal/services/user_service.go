package services

import (
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/repository"
)

// UserService 普通用户服务
type UserService struct {
	userRepo *repository.UserRepo
	config   *configs.Config
}

// NewUserService 创建普通用户服务
func NewUserService(userRepo *repository.UserRepo, config *configs.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   config,
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
		return nil, errorx.WrapError(err, fmt.Sprintf("检查用户名 %s 是否存在失败", req.Username))
	}
	if isExist {
		existsErr := errorx.Errorf(errorx.ErrUserExists, "用户名 %s 已存在", req.Username)
		return nil, errorx.WrapError(existsErr, "")
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errorx.WrapError(err, "密码加密失败")
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
		return nil, errorx.WrapError(err, fmt.Sprintf("创建用户 %s 失败", req.Username))
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		// 判断是否是用户不存在错误
		if errorx.IsAppErr(err) && err.(*errorx.AppError).GetErrorCode() == errorx.ErrorUserNotFoundCode {
			// 保持原始错误码，但添加业务上下文
			return nil, err
		}
		// 其他错误添加业务上下文
		return nil, errorx.WrapError(err, fmt.Sprintf("用户名 %s 登录失败", username))
	}

	// 检查用户是否启用
	if !user.Enabled {
		disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "用户 %s 已被禁用", username)
		return nil, errorx.WrapError(disabledErr, "")
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		passwordErr := errorx.Errorf(errorx.ErrUserPasswordError, "用户 %s 的密码错误", username)
		return nil, errorx.WrapError(passwordErr, "")
	}

	// 更新最后登录时间和IP
	fields := map[string]any{
		"last_login": time.Now(),
		"last_ip":    ip,
	}
	if err := s.userRepo.UpdateFields(user.ID, fields); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("更新用户 %s 的登录信息失败", username))
	}

	// 获取配置
	cfg := s.config

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
	fields := make(map[string]any, len(data))
	for k, v := range data {
		fields[k] = v
	}
	if err := s.userRepo.UpdateFields(int64(id), fields); err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新用户ID %d 的信息失败", id))
	}
	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id int64, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errorx.WrapError(err, fmt.Sprintf("获取用户ID %d 失败", id))
	}

	// 验证旧密码
	if !s.VerifyPassword(oldPassword, user.Password) {
		passwordErr := errorx.Errorf(errorx.ErrUserPasswordError, "原密码错误")
		return errorx.WrapError(passwordErr, "")
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return errorx.WrapError(err, "密码加密失败")
	}

	// 更新密码
	fields := map[string]any{
		"password": hashedPassword,
	}
	if err := s.userRepo.UpdateFields(id, fields); err != nil {
		return errorx.WrapError(err, fmt.Sprintf("更新用户ID %d 的密码失败", id))
	}
	return nil
}

// 该函数已移动到 controller/user.go 中的 UserInfo 方法
