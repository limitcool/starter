package dto

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 登录请求（别名，保持兼容性）
type LoginRequest = UserLoginRequest

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Mobile   string `json:"mobile"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
}

// UserChangePasswordRequest 修改密码请求
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	AccessToken       string   `json:"access_token"`              // 访问令牌
	RefreshToken      string   `json:"refresh_token"`             // 刷新令牌
	TokenType         string   `json:"token_type"`                // 令牌类型
	ExpiresIn         int64    `json:"expires_in"`                // 过期时间（秒）
	ExpireTime        int64    `json:"expire_time"`               // 过期时间戳
	RefreshExpiresIn  int64    `json:"refresh_expires_in"`        // 刷新令牌过期时间（秒）
	RefreshExpireTime int64    `json:"refresh_expire_time"`       // 刷新令牌过期时间戳
	Scope             string   `json:"scope,omitempty"`           // 权限范围
	UserID            int64    `json:"user_id"`                   // 用户ID
	Username          string   `json:"username"`                  // 用户名
	Email             string   `json:"email,omitempty"`           // 邮箱
	Mobile            string   `json:"mobile,omitempty"`          // 手机号
	Nickname          string   `json:"nickname,omitempty"`        // 昵称
	Avatar            string   `json:"avatar,omitempty"`          // 头像
	Roles             []string `json:"roles,omitempty"`           // 角色列表
	Permissions       []string `json:"permissions,omitempty"`     // 权限列表
	IsAdmin           bool     `json:"is_admin"`                  // 是否管理员
	LastLoginTime     string   `json:"last_login_time,omitempty"` // 上次登录时间
	LastLoginIP       string   `json:"last_login_ip,omitempty"`   // 上次登录IP
}

// LoginResponse 登录响应（别名，保持兼容性）
type LoginResponse = UserLoginResponse

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
