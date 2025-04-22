package services

import (
	"context"
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// SimpleAdminUserService 简单模式下的管理员用户服务
type SimpleAdminUserService struct {
	userRepo    *repository.UserRepo
	config      *configs.Config
	authService *AuthService
}

// NewSimpleAdminUserService 创建简单模式下的管理员用户服务
func NewSimpleAdminUserService(params ServiceParams, authService *AuthService) *SimpleAdminUserService {
	service := &SimpleAdminUserService{
		userRepo:    params.UserRepo,
		config:      params.Config,
		authService: authService,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "SimpleAdminUserService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "SimpleAdminUserService stopped")
			return nil
		},
	})

	return service
}

// VerifyPassword 验证用户密码
func (s *SimpleAdminUserService) VerifyPassword(password, hashedPassword string) bool {
	// 使用通用的 VerifyPassword 函数
	return VerifyPassword(password, hashedPassword)
}

// Login 管理员登录
func (s *SimpleAdminUserService) Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("用户名 %s 不存在", username))
	}

	// 检查是否是管理员
	if !user.IsAdmin {
		permErr := errorx.Errorf(errorx.ErrAccessDenied, "用户 %s 不是管理员", username)
		return nil, errorx.WrapError(permErr, "")
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
	if err := s.userRepo.UpdateFields(ctx, user.ID, fields); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("更新用户 %s 的登录信息失败", username))
	}

	// 简单模式下管理员角色固定为admin
	roles := []string{"admin"}

	// 生成令牌
	return s.authService.GenerateTokensWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, roles)
}

// RefreshToken 刷新令牌
func (s *SimpleAdminUserService) RefreshToken(ctx context.Context, refreshToken string) (*v1.LoginResponse, error) {
	// 验证刷新令牌
	claims, err := s.authService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errorx.WrapError(err, "刷新令牌无效")
	}

	// 检查用户类型
	if claims.UserType != enum.UserTypeAdminUser {
		typeErr := errorx.Errorf(errorx.ErrAccessDenied, "用户类型 %s 不是管理员", claims.UserType)
		return nil, errorx.WrapError(typeErr, "")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户ID %d 失败", claims.UserID))
	}

	// 检查是否是管理员
	if !user.IsAdmin {
		permErr := errorx.Errorf(errorx.ErrAccessDenied, "用户ID %d 不是管理员", claims.UserID)
		return nil, errorx.WrapError(permErr, "")
	}

	// 检查用户是否启用
	if !user.Enabled {
		disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "用户ID %d 已被禁用", claims.UserID)
		return nil, errorx.WrapError(disabledErr, "")
	}

	// 简单模式下管理员角色固定为admin
	roles := []string{"admin"}

	// 生成新的访问令牌
	accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, roles)
	if err != nil {
		return nil, err
	}
	return s.authService.CreateLoginResponse(accessToken, refreshToken), nil
}

// GetUserInfo 获取管理员用户信息
func (s *SimpleAdminUserService) GetUserInfo(ctx context.Context, id int64) (interface{}, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取用户ID %d 失败", id))
	}

	// 检查是否是管理员
	if !user.IsAdmin {
		permErr := errorx.Errorf(errorx.ErrAccessDenied, "用户ID %d 不是管理员", id)
		return nil, errorx.WrapError(permErr, "")
	}

	// 简单模式下管理员角色固定为admin
	user.RoleCodes = []string{"admin"}

	return user, nil
}
