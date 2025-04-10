package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"gorm.io/gorm"
)

// UserService 普通用户服务
type UserService struct {
}

// NewUserService 创建普通用户服务
func NewUserService() *UserService {
	return &UserService{}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := sqldb.Instance().DB().First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyPassword 验证用户密码
func (s *UserService) VerifyPassword(password, hashedPassword string) bool {
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
func (s *UserService) Register(req RegisterRequest) (*model.User, error) {
	isExist, err := model.NewUser().IsExist()
	if err != nil {
		return nil, errorx.ErrDatabaseQueryError.WithError(err)
	}
	if isExist {
		return nil, errorx.ErrUserAlreadyExists
	}

	// 哈希密码
	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errorx.ErrInternal.WithError(err)
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
		Birthday:   &req.Birthday,
		Address:    req.Address,
		RegisterIP: req.RegisterIP,
	}

	if err := user.Create(); err != nil {
		return nil, errorx.ErrDatabaseInsertError.WithError(err)
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string, ip string) (*LoginResponse, error) {
	// 获取用户
	user, err := model.NewUser().GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// 检查用户是否启用
	if !user.Enabled {
		return nil, errorx.ErrUserDisabled
	}

	// 验证密码
	if !s.VerifyPassword(password, user.Password) {
		return nil, errorx.ErrUserPasswordError.WithMsg("密码错误")
	}

	// 更新最后登录时间和IP
	sqldb.Instance().DB().Model(user).Updates(map[string]interface{}{
		"last_login": time.Now(),
		"last_ip":    ip,
	})

	// 获取配置
	cfg := core.Instance().Config()

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeUser.String(),    // 普通用户
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  enum.UserTypeUser.String(),     // 普通用户
		TokenType: enum.TokenTypeRefresh.String(), // 刷新令牌
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

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, data map[string]interface{}) error {
	// 不允许更新的字段
	delete(data, "id")
	delete(data, "username")
	delete(data, "password")
	delete(data, "created_at")
	delete(data, "deleted_at")

	// 更新用户信息
	return sqldb.Instance().DB().Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id int64, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(oldPassword, user.Password) {
		return errorx.ErrUserPasswordError.WithMsg("原密码错误")
	}

	// 哈希新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	return sqldb.Instance().DB().Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func GetUserInfo(ctx *gin.Context) {
	userId, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, errorx.ErrUserNoLogin)
		return
	}

	var user model.User
	if err := sqldb.Instance().DB().First(&user, userId).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	// 加载头像文件信息
	if user.AvatarFileID != 0 {
		var avatarFile model.File
		if err := sqldb.Instance().DB().First(&avatarFile, user.AvatarFileID).Error; err == nil {
			user.AvatarURL = avatarFile.URL
		}
	}

	response.Success(ctx, user)
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
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 2. 验证用户名是否已存在
	var count int64
	if err := sqldb.Instance().DB().Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	if count > 0 {
		response.Error(ctx, errorx.ErrUserAlreadyExists)
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

	if err := sqldb.Instance().DB().Create(&user).Error; err != nil {

		response.Error(ctx, errorx.ErrDatabaseInsertError)
		return
	}

	// 4. 返回成功
	type RegisterResult struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	}

	result := RegisterResult{
		ID:       user.ID,
		Username: user.Username,
	}

	response.Success(ctx, result)
}

// 用户登录
func UserLogin(ctx *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 2. 查询用户
	var user model.User
	if err := sqldb.Instance().DB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(ctx, errorx.ErrUserNameOrPasswordError)
		} else {
			response.Error(ctx, errorx.ErrDatabaseQueryError)
		}
		return
	}

	// 3. 验证密码
	// 注意：实际应用中应该验证加密后的密码
	if user.Password != req.Password {
		response.Error(ctx, errorx.ErrUserNameOrPasswordError)
		return
	}

	// 检查用户状态
	if !user.Enabled {
		response.Error(ctx, errorx.ErrUserDisabled)
		return
	}

	// 4. 更新登录信息
	now := time.Now()
	user.LastLogin = &now
	user.LastIP = ctx.ClientIP()
	if err := sqldb.Instance().DB().Save(&user).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	// 5. 生成token
	// 注意：实际应用中应该使用JWT生成token
	token := "mock_token_" + req.Username

	// 加载头像文件信息
	if user.AvatarFileID != 0 {
		var avatarFile model.File
		if err := sqldb.Instance().DB().First(&avatarFile, user.AvatarFileID).Error; err == nil {
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

	response.Success(ctx, result)
}

// 修改密码
func ChangePassword(ctx *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	// 2. 获取当前用户
	userId, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, errorx.ErrUserNoLogin)
		return
	}

	// 3. 查询用户
	var user model.User
	if err := sqldb.Instance().DB().First(&user, userId).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	// 4. 验证旧密码
	// 注意：实际应用中应该验证加密后的密码
	if user.Password != req.OldPassword {
		response.Error(ctx, errorx.ErrUserPasswordError)
		return
	}

	// 5. 更新密码
	// 注意：实际应用中应该对新密码进行加密
	hashedPassword, _ := crypto.HashPassword(req.NewPassword)
	if err := sqldb.Instance().DB().Model(&user).Update("password", hashedPassword).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
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
	db := sqldb.Instance().DB().Model(&model.User{})

	// 如果有关键字，添加搜索条件
	if keyword != "" {
		db.Where("username LIKE ? OR nickname LIKE ? OR mobile LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	var total int64
	if err := sqldb.Instance().DB().Count(&total).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	// 分页查询
	var pageNum, pageSizeNum int
	_, err1 := fmt.Sscanf(page, "%d", &pageNum)
	_, err2 := fmt.Sscanf(pageSize, "%d", &pageSizeNum)
	if err1 != nil || err2 != nil || pageNum <= 0 || pageSizeNum <= 0 {
		response.Error(ctx, errorx.ErrInvalidParams)
		return
	}

	if err := sqldb.Instance().DB().Offset((pageNum - 1) * pageSizeNum).Limit(pageSizeNum).Find(&users).Error; err != nil {
		response.Error(ctx, errorx.ErrDatabaseQueryError)
		return
	}

	// 获取用户头像URL
	for i := range users {
		if users[i].AvatarFileID != 0 {
			var avatarFile model.File
			if err := sqldb.Instance().DB().First(&avatarFile, users[i].AvatarFileID).Error; err == nil {
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

	response.Success(ctx, result)
}
