package services

import (
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/casbin"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/repository"
)

// SysUserService 用户服务
type SysUserService struct {
	sysUserRepo   *repository.SysUserRepo
	userRepo      *repository.UserRepo
	roleService   *RoleService
	casbinService casbin.Service
	config        *configs.Config
}

// NewSysUserService 创建用户服务
func NewSysUserService(sysUserRepo *repository.SysUserRepo, userRepo *repository.UserRepo, roleService *RoleService, casbinService casbin.Service, config *configs.Config) *SysUserService {
	return &SysUserService{
		sysUserRepo:   sysUserRepo,
		userRepo:      userRepo,
		roleService:   roleService,
		casbinService: casbinService,
		config:        config,
	}
}

// VerifyPassword 验证用户密码
func (s *SysUserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// Login 用户登录
func (s *SysUserService) Login(username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := s.sysUserRepo.GetByUsername(username)
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
	if err := s.sysUserRepo.UpdateFields(user.ID, fields); err != nil {
		return nil, errorx.WrapError(err, fmt.Sprintf("更新用户 %s 的登录信息失败", username))
	}

	// 获取配置
	cfg := s.config

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeSysUser,          // 系统用户
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
		Roles:     user.RoleCodes,                // 角色编码
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeSysUser,           // 系统用户
		TokenType: enum.TokenTypeRefresh.String(), // 刷新令牌
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, errorx.WrapError(err, "生成访问令牌失败")
	}

	refreshToken, err := jwtpkg.GenerateTokenWithCustomClaims(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, errorx.WrapError(err, "生成刷新令牌失败")
	}

	return &v1.LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpiresIn:         cfg.JwtAuth.AccessExpire,
		ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
		RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
		RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *SysUserService) RefreshToken(refreshToken string) (*v1.LoginResponse, error) {
	// 获取配置
	cfg := s.config

	// 验证刷新令牌
	claims, err := jwtpkg.ParseTokenWithCustomClaims(refreshToken, cfg.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, errorx.WrapError(err, "验证刷新令牌失败")
	}

	// 检查令牌类型
	if claims.TokenType != enum.TokenTypeRefresh.String() {
		tokenErr := errorx.Errorf(errorx.ErrUserTokenError, "无效的令牌类型")
		return nil, errorx.WrapError(tokenErr, "")
	}

	// 获取用户类型
	userType := claims.UserType

	// 定义要返回的响应对象
	var loginResponse *v1.LoginResponse

	// 根据用户类型不同，查询不同的表
	switch userType {
	case enum.UserTypeSysUser:
		// 系统用户 - 查询系统用户表
		user, err := s.sysUserRepo.GetByID(claims.UserID)
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
		accessClaims := &jwtpkg.CustomClaims{
			UserID:    user.ID,
			Username:  user.Username,
			UserType:  enum.UserTypeSysUser,          // 系统用户
			TokenType: enum.TokenTypeAccess.String(), // 访问令牌
			Roles:     user.RoleCodes,                // 角色编码
		}

		accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
		if err != nil {
			return nil, errorx.WrapError(err, "生成系统用户访问令牌失败")
		}

		loginResponse = &v1.LoginResponse{
			AccessToken:       accessToken,
			RefreshToken:      refreshToken, // 保持原有的刷新令牌
			ExpiresIn:         cfg.JwtAuth.AccessExpire,
			ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
			RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
			RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
		}

	case enum.UserTypeUser:
		// 普通用户 - 查询普通用户表
		user, err := s.userRepo.GetByID(claims.UserID)
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
		accessClaims := &jwtpkg.CustomClaims{
			UserID:    user.ID,
			Username:  user.Username,
			UserType:  enum.UserTypeUser,             // 普通用户
			TokenType: enum.TokenTypeAccess.String(), // 访问令牌
			Roles:     []string{"user"},              // 普通用户默认角色
		}

		accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
		if err != nil {
			return nil, errorx.WrapError(err, "生成普通用户访问令牌失败")
		}

		loginResponse = &v1.LoginResponse{
			AccessToken:       accessToken,
			RefreshToken:      refreshToken, // 保持原有的刷新令牌
			ExpiresIn:         cfg.JwtAuth.AccessExpire,
			ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
			RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
			RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
		}

	default:
		tokenErr := errorx.Errorf(errorx.ErrUserTokenError, "无效的用户类型: %s", userType)
		return nil, errorx.WrapError(tokenErr, "")
	}

	return loginResponse, nil
}
