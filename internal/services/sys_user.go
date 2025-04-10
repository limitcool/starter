package services

import (
	"errors"
	"time"

	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
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

// GetUserByID 根据ID获取用户
func (s *SysUserService) GetUserByID(id int64) (*model.SysUser, error) {
	var user model.SysUser
	err := sqldb.Instance().DB().First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// 获取用户的角色
	if err := sqldb.Instance().DB().Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *SysUserService) GetUserByUsername(username string) (*model.SysUser, error) {
	var user model.SysUser
	err := sqldb.Instance().DB().Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		// 直接返回错误，错误会自动捕获堆栈
		return nil, err
	}

	// // 获取用户的角色
	// if err := sqldb.Instance().DB().Model(&user).Association("Roles").Find(&user.Roles); err != nil {
	// 	// 直接返回错误，错误会自动捕获堆栈
	// 	return nil, err
	// }

	// // 提取角色编码
	// for _, role := range user.Roles {
	// 	user.RoleCodes = append(user.RoleCodes, role.Code)
	// }

	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *SysUserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	ExpiresIn         int64  `json:"expires_in"`          // 访问令牌有效期（秒）
	ExpireTime        int64  `json:"expire_time"`         // 访问令牌过期时间戳
	RefreshExpiresIn  int64  `json:"refresh_expires_in"`  // 刷新令牌有效期（秒）
	RefreshExpireTime int64  `json:"refresh_expire_time"` // 刷新令牌过期时间戳
}

// Login 用户登录
func (s *SysUserService) Login(username, password string, ip string) (*LoginResponse, error) {
	// 获取用户
	user, err := s.GetUserByUsername(username)
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
		UserType:  enum.UserTypeSysUser.String(), // 系统用户
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
		Roles:     user.RoleCodes,                // 角色编码
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeSysUser.String(),  // 系统用户
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

	return &LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpiresIn:         cfg.JwtAuth.AccessExpire,
		ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
		RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
		RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *SysUserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
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

	// 获取用户信息
	user, err := s.GetUserByID(claims.UserID)
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
		UserType:  enum.UserTypeSysUser.String(), // 系统用户
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
		Roles:     user.RoleCodes,                // 角色编码
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		// AppError的WithError会自动捕获堆栈
		return nil, errorx.ErrInternal.WithError(err)
	}

	return &LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken, // 保持原有的刷新令牌
		ExpiresIn:         cfg.JwtAuth.AccessExpire,
		ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
		RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
		RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}, nil
}
