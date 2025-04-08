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

// NormalUserService 普通用户服务
type NormalUserService struct {
	db *gorm.DB
}

// NewNormalUserService 创建普通用户服务
func NewNormalUserService(db *gorm.DB) *NormalUserService {
	return &NormalUserService{
		db: db,
	}
}

// GetUserByID 根据ID获取用户
func (s *NormalUserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.NewErrCodeMsg(code.UserNotFound, "用户不存在")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *NormalUserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, code.NewErrCodeMsg(code.UserNotFound, "用户不存在")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *NormalUserService) VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string    `json:"username" binding:"required"`
	Password   string    `json:"password" binding:"required"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Mobile     string    `json:"mobile"`
	Gender     string    `json:"gender"`
	Birthday   time.Time `json:"birthday"`
	Address    string    `json:"address"`
	RegisterIP string    `json:"register_ip"`
}

// Register 用户注册
func (s *NormalUserService) Register(req RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, code.NewErrCodeMsg(code.UserAlreadyExists, "用户名已存在")
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
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
		Birthday:   req.Birthday,
		Address:    req.Address,
		RegisterIP: req.RegisterIP,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *NormalUserService) Login(username, password string, ip string) (*LoginResponse, error) {
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

	// 更新最后登录时间和IP
	s.db.Model(user).Updates(map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	})

	// 生成访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": "user", // 普通用户
		"type":      "access_token",
		"exp":       time.Now().Add(time.Duration(global.Config.JwtAuth.AccessExpire) * time.Second).Unix(),
	}

	// 生成刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": "user", // 普通用户
		"type":      "refresh_token",
		"exp":       time.Now().Add(time.Duration(global.Config.JwtAuth.RefreshExpire) * time.Second).Unix(),
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

// UpdateUser 更新用户信息
func (s *NormalUserService) UpdateUser(id uint, data map[string]interface{}) error {
	// 不允许更新的字段
	delete(data, "id")
	delete(data, "username")
	delete(data, "password")
	delete(data, "created_at")
	delete(data, "deleted_at")

	// 更新用户信息
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

// ChangePassword 修改密码
func (s *NormalUserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(oldPassword, user.Password) {
		return code.NewErrCodeMsg(code.UserPasswordError, "原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}
