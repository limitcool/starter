package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
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
		return nil, errorx.NewErrCodeMsg(errorx.UserNotFound, "用户不存在")
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
		return nil, errorx.NewErrCodeMsg(errorx.UserNotFound, "用户不存在")
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
		return nil, errorx.NewErrCodeMsg(errorx.UserAlreadyExists, "用户名已存在")
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

// 获取配置，优先使用ServiceManager
func (s *NormalUserService) getConfig() *configs.Config {
	if serviceInstance != nil {
		return serviceInstance.GetConfig()
	}
	panic("配置未初始化")
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
		return nil, errorx.NewErrCodeMsg(errorx.UserDisabled, "用户已禁用")
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, errorx.NewErrCodeMsg(errorx.UserPasswordError, "密码错误")
	}

	// 更新最后登录时间和IP
	s.db.Model(user).Updates(map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	})

	// 获取配置
	cfg := s.getConfig()

	// 生成访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": "user", // 普通用户
		"type":      "access_token",
		"exp":       time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
	}

	// 生成刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": "user", // 普通用户
		"type":      "refresh_token",
		"exp":       time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}

	accessToken, err := jwtpkg.GenerateToken(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := jwtpkg.GenerateToken(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    cfg.JwtAuth.AccessExpire,
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
		return errorx.NewErrCodeMsg(errorx.UserPasswordError, "原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func GetUserInfo(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		response.Fail(ctx, errorx.UserNoLogin, "")
		return
	}

	var user model.User
	if err := model.GetDB().First(&user, userId).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 加载头像文件信息
	if user.Avatar != "" {
		var avatarFile model.File
		if err := model.GetDB().First(&avatarFile, user.Avatar).Error; err == nil {
			user.AvatarURL = avatarFile.URL
		}
	}

	response.Success[model.User](ctx, user)
}

// 用户注册
func UserRegister(ctx *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Mobile   string `json:"mobile"`
		Email    string `json:"email"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, errorx.InvalidParams, err.Error())
		return
	}

	// 2. 验证用户名是否已存在
	var count int64
	if err := model.GetDB().Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	if count > 0 {
		response.Fail(ctx, errorx.UserAlreadyExists, "")
		return
	}

	// 3. 创建新用户
	user := model.User{
		Username:   req.Username,
		Password:   req.Password, // 注意：实际应用中应该对密码进行加密
		Nickname:   req.Nickname,
		Mobile:     req.Mobile,
		Email:      req.Email,
		Enabled:    true,
		RegisterIP: ctx.ClientIP(),
	}

	if err := model.GetDB().Create(&user).Error; err != nil {
		if errorx.IsDuplicate(err) {
			response.Fail(ctx, errorx.UserAlreadyExists, "")
			return
		}
		response.Fail(ctx, errorx.DatabaseInsertError, "")
		return
	}

	// 4. 返回成功
	type RegisterResult struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	result := RegisterResult{
		ID:       user.ID,
		Username: user.Username,
	}

	response.Success[RegisterResult](ctx, result)
}

// 用户登录
func UserLogin(ctx *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, errorx.InvalidParams, "")
		return
	}

	// 2. 查询用户
	var user model.User
	if err := model.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(ctx, errorx.UserNameOrPasswordError, "")
		} else {
			response.Fail(ctx, errorx.DatabaseQueryError, "")
		}
		return
	}

	// 3. 验证密码
	// 注意：实际应用中应该验证加密后的密码
	if user.Password != req.Password {
		response.Fail(ctx, errorx.UserNameOrPasswordError, "")
		return
	}

	// 检查用户状态
	if !user.Enabled {
		response.Fail(ctx, errorx.UserDisabled, "")
		return
	}

	// 4. 更新登录信息
	user.LastLogin = time.Now()
	user.LastIP = ctx.ClientIP()
	if err := model.GetDB().Save(&user).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 5. 生成token
	// 注意：实际应用中应该使用JWT生成token
	token := "mock_token_" + req.Username

	// 加载头像文件信息
	if user.Avatar != "" {
		var avatarFile model.File
		if err := model.GetDB().First(&avatarFile, user.Avatar).Error; err == nil {
			user.AvatarURL = avatarFile.URL
		}
	}

	// 6. 返回成功
	type LoginResult struct {
		Token    string     `json:"token"`
		UserInfo model.User `json:"user_info"`
	}

	result := LoginResult{
		Token:    token,
		UserInfo: user,
	}

	response.Success[LoginResult](ctx, result)
}

// 修改密码
func ChangePassword(ctx *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, errorx.InvalidParams, "")
		return
	}

	// 2. 获取当前用户
	userId, exists := ctx.Get("userID")
	if !exists {
		response.Fail(ctx, errorx.UserNoLogin, "")
		return
	}

	// 3. 查询用户
	var user model.User
	if err := model.GetDB().First(&user, userId).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 4. 验证旧密码
	// 注意：实际应用中应该验证加密后的密码
	if user.Password != req.OldPassword {
		response.Fail(ctx, errorx.UserPasswordError, "")
		return
	}

	// 5. 更新密码
	// 注意：实际应用中应该对新密码进行加密
	hashedPassword, _ := crypto.HashPassword(req.NewPassword)
	if err := model.GetDB().Model(&user).Update("password", hashedPassword).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 6. 返回成功
	response.Success[any](ctx, nil)
}

// 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 1. 解析分页参数
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	keyword := ctx.DefaultQuery("keyword", "")

	// 2. 查询用户
	var users []model.User
	db := model.GetDB().Model(&model.User{})

	// 如果有关键字，添加搜索条件
	if keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR mobile LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 分页查询
	var pageNum, pageSizeNum int
	_, err1 := fmt.Sscanf(page, "%d", &pageNum)
	_, err2 := fmt.Sscanf(pageSize, "%d", &pageSizeNum)
	if err1 != nil || err2 != nil || pageNum <= 0 || pageSizeNum <= 0 {
		response.Fail(ctx, errorx.InvalidParams, "分页参数无效")
		return
	}

	if err := db.Offset((pageNum - 1) * pageSizeNum).Limit(pageSizeNum).Find(&users).Error; err != nil {
		response.Fail(ctx, errorx.DatabaseQueryError, "")
		return
	}

	// 获取用户头像URL
	for i := range users {
		if users[i].Avatar != "" {
			var avatarFile model.File
			if err := model.GetDB().First(&avatarFile, users[i].Avatar).Error; err == nil {
				users[i].AvatarURL = avatarFile.URL
			}
		}
	}

	// 3. 返回结果
	type ResponseData struct {
		Total int64        `json:"total"`
		List  []model.User `json:"list"`
	}

	result := ResponseData{
		Total: total,
		List:  users,
	}

	response.Success[ResponseData](ctx, result)
}
