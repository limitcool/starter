package services

import (
	"time"

	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/storage/sqldb"
)

// SysUserService 用户服务
type SysUserService struct {
	roleService   *RoleService
	casbinService *CasbinService
}

// NewSysUserService 创建用户服务
func NewSysUserService() *SysUserService {
	// 直接创建依赖服务
	return &SysUserService{
		roleService:   NewRoleService(),
		casbinService: NewCasbinService(),
	}
}

// VerifyPassword 验证用户密码
func (s *SysUserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// Login 用户登录
func (s *SysUserService) Login(username, password string, ip string) (*v1.LoginResponse, error) {
	// 获取用户
	user, err := model.NewSysUser().GetUserByUsername(username)
	if err != nil {
		// AppError的WithError会自动捕获堆栈
		return nil, errorx.ErrUserNotFound.WithError(err)
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, errorx.ErrUserDisabled
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, errorx.ErrUserPasswordError
	}
	db := sqldb.Instance().DB()
	// 更新最后登录时间和IP
	if err := db.Model(user).Updates(map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	}).Error; err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// 获取配置
	cfg := core.Instance().Config()

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
		// AppError的WithError会自动捕获堆栈
		return nil, errorx.ErrInternal.WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateTokenWithCustomClaims(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		// AppError的WithError会自动捕获堆栈
		return nil, errorx.ErrInternal.WithError(err)
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
	cfg := core.Instance().Config()

	// 验证刷新令牌
	claims, err := jwtpkg.ParseTokenWithCustomClaims(refreshToken, cfg.JwtAuth.RefreshSecret)
	if err != nil {
		// AppError的WithError会自动捕获堆栈
		return nil, errorx.ErrUserTokenError.WithError(err)
	}

	// 检查令牌类型
	if claims.TokenType != enum.TokenTypeRefresh.String() {
		return nil, errorx.ErrUserTokenError.WithMsg("无效的令牌类型")
	}

	// 获取用户类型
	userType := claims.UserType

	// 定义要返回的响应对象
	var loginResponse *v1.LoginResponse

	// 根据用户类型不同，查询不同的表
	switch userType {
	case enum.UserTypeSysUser:
		// 系统用户 - 查询系统用户表
		user, err := model.NewSysUser().GetUserByID(claims.UserID)
		if err != nil {
			if errorx.IsAppErr(err) {
				return nil, err
			}
			// AppError的WithError会自动捕获堆栈
			return nil, errorx.ErrUserNotFound.WithError(err)
		}

		// 检查用户是否启用
		if !user.Enabled {
			return nil, errorx.ErrUserDisabled
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
			// AppError的WithError会自动捕获堆栈
			return nil, errorx.ErrInternal.WithError(err)
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
		user, err := model.NewUser().GetUserByID(claims.UserID)
		if err != nil {
			if errorx.IsAppErr(err) {
				return nil, err
			}
			return nil, errorx.ErrUserNotFound.WithError(err)
		}

		// 检查用户状态（如果普通用户有状态字段）
		if !user.Enabled {
			return nil, errorx.ErrUserDisabled
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
			return nil, errorx.ErrInternal.WithError(err)
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
		return nil, errorx.ErrUserTokenError.WithMsg("无效的用户类型")
	}

	return loginResponse, nil
}
