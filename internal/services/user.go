package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/pkg/code"
	"github.com/limitcool/starter/pkg/crypto"
	jwtpkg "github.com/limitcool/starter/pkg/jwt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db            *gorm.DB
	roleService   *RoleService
	casbinService *CasbinService
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db:            db,
		roleService:   NewRoleService(db),
		casbinService: NewCasbinService(db),
	}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.NewErrCodeMsg(code.UserNotFound, "用户不存在")
	}
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := s.db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.NewErrCodeMsg(code.UserNotFound, "用户不存在")
	}
	if err != nil {
		return nil, err
	}

	// 获取用户的角色
	if err := s.db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}

	// 提取角色编码
	for _, role := range user.Roles {
		user.RoleCodes = append(user.RoleCodes, role.Code)
	}

	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *UserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*LoginResponse, error) {
	// 获取用户
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, code.NewErrCodeMsg(code.UserDisabled, "用户已禁用")
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, code.NewErrCodeMsg(code.UserPasswordError, "密码错误")
	}

	// 更新最后登录时间
	s.db.Model(user).Update("last_login", time.Now())

	// 生成访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role_codes": user.RoleCodes,
		"type":       "access_token",
		"exp":        time.Now().Add(time.Duration(global.Config.JwtAuth.AccessExpire) * time.Second).Unix(),
	}

	// 生成刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"type":     "refresh_token",
		"exp":      time.Now().Add(time.Duration(global.Config.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}

	accessToken, err := jwtpkg.GenerateToken(accessClaims, global.Config.JwtAuth.AccessSecret, time.Duration(global.Config.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := jwtpkg.GenerateToken(refreshClaims, global.Config.JwtAuth.RefreshSecret, time.Duration(global.Config.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    global.Config.JwtAuth.AccessExpire,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *UserService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// 验证刷新令牌
	tokenClaims, err := jwtpkg.ParseToken(refreshToken, global.Config.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, code.NewErrCodeMsg(code.UserTokenError, "无效的刷新令牌")
	}

	// 检查令牌类型
	claims := *tokenClaims
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh_token" {
		return nil, code.NewErrCodeMsg(code.UserTokenError, "无效的令牌类型")
	}

	// 获取用户信息
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, code.NewErrCodeMsg(code.UserTokenError, "无效的用户信息")
	}

	user, err := s.GetUserByID(uint(userID))
	if err != nil {
		return nil, err
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, code.NewErrCodeMsg(code.UserDisabled, "用户已禁用")
	}

	// 生成新的访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role_codes": user.RoleCodes,
		"type":       "access_token",
		"exp":        time.Now().Add(time.Duration(global.Config.JwtAuth.AccessExpire) * time.Second).Unix(),
	}

	accessToken, err := jwtpkg.GenerateToken(accessClaims, global.Config.JwtAuth.AccessSecret, time.Duration(global.Config.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // 保持原有的刷新令牌
		ExpiresIn:    global.Config.JwtAuth.AccessExpire,
	}, nil
}
