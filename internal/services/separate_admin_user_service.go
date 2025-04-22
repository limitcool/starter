package services

import (
	"context"
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/repository"
	"go.uber.org/fx"
)

// SeparateAdminUserService 分离模式下的管理员用户服务
type SeparateAdminUserService struct {
	adminUserRepo *repository.AdminUserRepo
	roleRepo      *repository.RoleRepo
	casbinService casbin.Service
	config        *configs.Config
	authService   *AuthService
}

// NewSeparateAdminUserService 创建分离模式下的管理员用户服务
func NewSeparateAdminUserService(params ServiceParams, casbinService casbin.Service, authService *AuthService) *SeparateAdminUserService {
	service := &SeparateAdminUserService{
		adminUserRepo: params.AdminUserRepo,
		roleRepo:      params.RoleRepo,
		casbinService: casbinService,
		config:        params.Config,
		authService:   authService,
	}

	// 注册生命周期钩子
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.InfoContext(ctx, "SeparateAdminUserService initialized")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.InfoContext(ctx, "SeparateAdminUserService stopped")
			return nil
		},
	})

	return service
}

// VerifyPassword 验证用户密码
func (s *SeparateAdminUserService) VerifyPassword(password, hashedPassword string) bool {
	// 使用通用的 VerifyPassword 函数
	return VerifyPassword(password, hashedPassword)
}

// Login 管理员登录
func (s *SeparateAdminUserService) Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取管理员用户
	user, err := s.adminUserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("管理员用户名 %s 不存在", username))
	}

	// 检查用户是否启用
	if !user.Enabled {
		disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "管理员用户 %s 已被禁用", username)
		return nil, errorx.WrapError(disabledErr, "")
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		passwordErr := errorx.Errorf(errorx.ErrUserPasswordError, "管理员用户 %s 的密码错误", username)
		return nil, errorx.WrapError(passwordErr, "")
	}

	// 更新最后登录时间和IP
	fields := map[string]any{
		"last_login": time.Now(),
		"last_ip":    ip,
	}
	if err := s.adminUserRepo.UpdateFields(ctx, user.ID, fields); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("更新管理员用户 %s 的登录信息失败", username))
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetRoleCodesByAdminUserID(ctx, user.ID)
	if err != nil {
		// 如果获取角色失败，使用默认角色
		roles = []string{"admin"}
	}

	// 生成令牌
	return s.authService.GenerateTokensWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, roles)
}

// RefreshToken 刷新令牌
func (s *SeparateAdminUserService) RefreshToken(ctx context.Context, refreshToken string) (*v1.LoginResponse, error) {
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

	// 获取管理员用户
	user, err := s.adminUserRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户ID %d 失败", claims.UserID))
	}

	// 检查用户是否启用
	if !user.Enabled {
		disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "管理员用户ID %d 已被禁用", claims.UserID)
		return nil, errorx.WrapError(disabledErr, "")
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetRoleCodesByAdminUserID(ctx, user.ID)
	if err != nil {
		// 如果获取角色失败，使用默认角色
		roles = []string{"admin"}
	}

	// 生成新的访问令牌
	accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, roles)
	if err != nil {
		return nil, err
	}
	return s.authService.CreateLoginResponse(accessToken, refreshToken), nil
}

// GetUserInfo 获取管理员用户信息
func (s *SeparateAdminUserService) GetUserInfo(ctx context.Context, id int64) (interface{}, error) {
	// 获取管理员用户
	user, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户ID %d 失败", id))
	}

	// 获取用户角色
	roleCodes, err := s.roleRepo.GetRoleCodesByAdminUserID(ctx, user.ID)
	if err != nil {
		// 如果获取角色失败，使用空列表
		roleCodes = []string{}
	}
	user.RoleCodes = roleCodes

	return user, nil
}
