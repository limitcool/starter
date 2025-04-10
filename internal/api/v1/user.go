package v1

// UserRegisterRequest 用户注册请求参数
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
}

// UserChangePasswordRequest 用户修改密码请求参数
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
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
