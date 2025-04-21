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

// AdminUserService 管理员用户服务
type AdminUserService struct {
	adminUserRepo *repository.AdminUserRepo
	userRepo      *repository.UserRepo
	roleService   *RoleService
	casbinService casbin.Service
	config        *configs.Config
	authService   *AuthService
}

// NewAdminUserService 创建管理员用户服务
func NewAdminUserService(adminUserRepo *repository.AdminUserRepo, userRepo *repository.UserRepo, roleService *RoleService, casbinService casbin.Service, config *configs.Config) *AdminUserService {
	// 创建认证服务
	authService := NewAuthService(config)

	return &AdminUserService{
		adminUserRepo: adminUserRepo,
		userRepo:      userRepo,
		roleService:   roleService,
		casbinService: casbinService,
		config:        config,
		authService:   authService,
	}
}

// VerifyPassword 验证用户密码
func (s *AdminUserService) VerifyPassword(password, hashedPassword string) bool {
	// 使用通用的 VerifyPassword 函数
	return VerifyPassword(password, hashedPassword)
}

// Login 管理员用户登录
func (s *AdminUserService) Login(ctx context.Context, username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户模式
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 根据用户模式处理
	if userMode == enum.UserModeSimple {
		// 简单模式，使用普通用户表查询管理员用户
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

		// 获取用户角色
		roles, err := s.userRepo.GetUserRoles(ctx, user.ID, user.IsAdmin, s.config.Admin.UserMode)
		if err != nil {
			// 如果获取角色失败，使用默认角色
			roles = []string{"admin"}
		}

		// 生成令牌
		return s.authService.GenerateTokensWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, roles)
	} else {
		// 分离模式，使用管理员用户表
		user, err := s.adminUserRepo.GetByUsername(ctx, username)
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
		if err := s.adminUserRepo.UpdateFields(ctx, user.ID, fields); err != nil {
			return nil, errorx.WrapError(err, fmt.Sprintf("更新用户 %s 的登录信息失败", username))
		}

		// 使用认证服务生成令牌
		return s.authService.GenerateTokensWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, user.RoleCodes)
	}
}

// RefreshToken 刷新访问令牌
func (s *AdminUserService) RefreshToken(ctx context.Context, refreshToken string) (*v1.LoginResponse, error) {
	// 验证刷新令牌
	claims, err := s.authService.ValidateRefreshTokenWithContext(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户类型和用户模式
	userType := claims.UserType
	userMode := enum.GetUserMode(s.config.Admin.UserMode)

	// 根据用户类型和用户模式处理
	switch userType {
	case enum.UserTypeAdminUser, enum.UserTypeSysUser: // 支持旧版本的UserTypeSysUser
		// 如果是简单模式，从普通用户表中查询管理员用户
		if userMode == enum.UserModeSimple {
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

			// 获取用户角色
			roles, roleErr := s.userRepo.GetUserRoles(ctx, user.ID, user.IsAdmin, s.config.Admin.UserMode)
			if roleErr != nil {
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

		// 分离模式 - 查询管理员用户表
		user, err := s.adminUserRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			if errorx.IsAppErr(err) {
				return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户ID %d 失败", claims.UserID))
			}
			return nil, errorx.WrapError(err, fmt.Sprintf("获取管理员用户ID %d 失败", claims.UserID))
		}

		// 检查用户是否启用
		if !user.Enabled {
			disabledErr := errorx.Errorf(errorx.ErrUserDisabled, "管理员用户ID %d 已被禁用", claims.UserID)
			return nil, errorx.WrapError(disabledErr, "")
		}

		// 生成新的访问令牌
		accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, enum.UserTypeAdminUser, user.RoleCodes)
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

		// 获取用户角色
		roles, roleErr := s.userRepo.GetUserRoles(ctx, user.ID, user.IsAdmin, s.config.Admin.UserMode)
		if roleErr != nil {
			// 如果获取角色失败，使用默认角色
			if user.IsAdmin {
				roles = []string{"admin"}
			} else {
				roles = []string{"user"}
			}
		}

		// 判断用户类型
		userType := enum.UserTypeUser
		if user.IsAdmin {
			userType = enum.UserTypeAdminUser
		}

		// 生成新的访问令牌
		accessToken, err := s.authService.GenerateAccessTokenWithContext(ctx, uint(user.ID), user.Username, userType, roles)
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
