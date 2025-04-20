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
	"github.com/limitcool/starter/internal/repository"
)

// SysUserService 用户服务
type SysUserService struct {
	sysUserRepo   *repository.SysUserRepo
	userRepo      *repository.UserRepo
	roleService   *RoleService
	casbinService casbin.Service
	config        *configs.Config
	authService   *AuthService
}

// NewSysUserService 创建用户服务
func NewSysUserService(sysUserRepo *repository.SysUserRepo, userRepo *repository.UserRepo, roleService *RoleService, casbinService casbin.Service, config *configs.Config) *SysUserService {
	// 创建认证服务
	authService := NewAuthService(config)

	return &SysUserService{
		sysUserRepo:   sysUserRepo,
		userRepo:      userRepo,
		roleService:   roleService,
		casbinService: casbinService,
		config:        config,
		authService:   authService,
	}
}

// VerifyPassword 验证用户密码
func (s *SysUserService) VerifyPassword(password, hashedPassword string) bool {
	// 使用通用的 VerifyPassword 函数
	return VerifyPassword(password, hashedPassword)
}

// Login 用户登录
func (s *SysUserService) Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.sysUserRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("用户名 %s 不存在", username))
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
	if err := s.sysUserRepo.UpdateFields(ctx, user.ID, fields); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("更新用户 %s 的登录信息失败", username))
	}

	// 使用认证服务生成令牌
	return s.authService.GenerateTokensWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeSysUser, user.RoleCodes)
}

// RefreshToken 刷新访问令牌
func (s *SysUserService) RefreshToken(ctx context.Context, refreshToken string) (*v1.LoginResponse, error) {
	// 验证刷新令牌
	claims, err := s.authService.ValidateRefreshTokenWithContext(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户类型
	userType := claims.UserType

	// 根据用户类型不同，查询不同的表
	switch userType {
	case enum.UserTypeSysUser:
		// 系统用户 - 查询系统用户表
		user, err := s.sysUserRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			if errorx.IsAppErr(err) {
				return nil, errorx.WrapError(err, fmt.Sprintf("获取系统用户ID %d 失败", claims.UserID))
			}
			return nil, errorx.WrapError(err, fmt.Sprintf("获取系统用户ID %d 失败", claims.UserID))
		}

		// 检查用户是否启用
		if !user.Enabled {
			disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "系统用户ID %d 已被禁用", claims.UserID)
			return nil, errorx.WrapError(disabledErr, "")
		}

		// 生成新的访问令牌
		accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeSysUser, user.RoleCodes)
		if err != nil {
			return nil, err
		}

		// 创建登录响应
		return s.authService.CreateLoginResponse(accessToken, refreshToken), nil

	case enum.UserTypeUser:
		// 普通用户 - 查询普通用户表
		user, err := s.userRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			if errorx.IsAppErr(err) {
				return nil, errorx.WrapError(err, fmt.Sprintf("获取普通用户ID %d 失败", claims.UserID))
			}
			return nil, errorx.WrapError(err, fmt.Sprintf("获取普通用户ID %d 失败", claims.UserID))
		}

		// 检查用户状态
		if !user.Enabled {
			disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "普通用户ID %d 已被禁用", claims.UserID)
			return nil, errorx.WrapError(disabledErr, "")
		}

		// 生成新的访问令牌
		accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeUser, []string{"user"})
		if err != nil {
			return nil, err
		}

		// 创建登录响应
		return s.authService.CreateLoginResponse(accessToken, refreshToken), nil

	default:
		tokenErr := errorx.Errorf(errorx.ErrUserTokenError, "无效的用户类型: %s", userType)
		return nil, errorx.WrapError(tokenErr, "")
	}
}
