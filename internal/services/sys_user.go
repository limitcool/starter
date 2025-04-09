package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"gorm.io/gorm"
)

// SysUserService 用户服务
type SysUserService struct {
	roleService   *RoleService
	casbinService *CasbinService
}

// NewSysUserService 创建用户服务
func NewSysUserService() *SysUserService {

	// 使用ServiceManager获取依赖服务
	return &SysUserService{
		roleService:   serviceInstance.GetRoleService(),
		casbinService: serviceInstance.GetCasbinService(),
	}

}

// GetConfig 获取配置
func (s *SysUserService) getConfig() *configs.Config {
	if serviceInstance != nil {
		return serviceInstance.GetConfig()
	}
	panic("配置未初始化")
}

// GetUserByID 根据ID获取用户
func (s *SysUserService) GetUserByID(id uint) (*model.SysUser, error) {
	var user model.SysUser
	err := Instance().GetDB().First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := Instance().GetDB().Model(&user).Association("Roles").Find(&user.Roles); err != nil {
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
	err := Instance().GetDB().Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := Instance().GetDB().Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *SysUserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login 用户登录
func (s *SysUserService) Login(username, password string, ip string) (*LoginResponse, error) {
	// 获取用户
	user, err := s.GetUserByUsername(username)
	if err != nil {
		log.Error("登录失败", "err", err)
		return nil, errorx.ErrUserNotFound
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, errorx.ErrUserDisabled
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, errorx.ErrUserPasswordError
	}

	// 更新最后登录时间和IP
	Instance().GetDB().Model(user).Updates(map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	})

	// 获取配置
	cfg := s.getConfig()

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeSysUser.String(), // 系统用户
		TokenType: "access_token",                // 访问令牌
		Roles:     user.RoleCodes,                // 角色编码
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeSysUser.String(), // 系统用户
		TokenType: "refresh_token",               // 刷新令牌
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := jwtpkg.GenerateTokenWithCustomClaims(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    cfg.JwtAuth.AccessExpire,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *SysUserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// 获取配置
	cfg := s.getConfig()

	// 验证刷新令牌
	claims, err := jwtpkg.ParseTokenWithCustomClaims(refreshToken, cfg.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, errorx.ErrUserTokenError
	}

	// 检查令牌类型
	if claims.TokenType != enum.TokenTypeRefresh.String() {
		return nil, errorx.ErrUserTokenError.WithMsg("无效的令牌类型")
	}

	// 获取用户信息
	user, err := s.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // 保持原有的刷新令牌
		ExpiresIn:    cfg.JwtAuth.AccessExpire,
	}, nil
}
