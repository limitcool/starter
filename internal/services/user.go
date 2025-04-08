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

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	// 获取用户
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	// 检查用户是否启用
	if !user.Enabled {
		return "", code.NewErrCodeMsg(code.UserDisabled, "用户已禁用")
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return "", code.NewErrCodeMsg(code.UserPasswordError, "密码错误")
	}

	// 更新最后登录时间
	s.db.Model(user).Update("last_login", time.Now())

	// 生成JWT Token
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"username":   user.Username,
		"role_codes": user.RoleCodes,
	}

	expireDuration := time.Duration(global.Config.JwtAuth.AccessExpire) * time.Second
	token, err := jwtpkg.GenerateToken(claims, global.Config.JwtAuth.AccessSecret, expireDuration)
	if err != nil {
		return "", fmt.Errorf("生成token失败: %w", err)
	}

	return token, nil
}
