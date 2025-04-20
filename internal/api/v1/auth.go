package v1

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken       string   `json:"access_token"`              // 访问令牌
	RefreshToken      string   `json:"refresh_token"`             // 刷新令牌
	ExpiresIn         int64    `json:"expires_in"`                // 过期时间（秒）
	ExpireTime        int64    `json:"expire_time"`               // 过期时间戳
	RefreshExpiresIn  int64    `json:"refresh_expires_in"`        // 刷新令牌过期时间（秒）
	RefreshExpireTime int64    `json:"refresh_expire_time"`       // 刷新令牌过期时间戳
	TokenType         string   `json:"token_type"`                // 令牌类型
	Scope             string   `json:"scope,omitempty"`           // 权限范围
	UserID            int64    `json:"user_id"`                   // 用户ID
	Username          string   `json:"username"`                  // 用户名
	Nickname          string   `json:"nickname,omitempty"`        // 昵称
	Avatar            string   `json:"avatar,omitempty"`          // 头像
	Roles             []string `json:"roles,omitempty"`           // 角色列表
	Permissions       []string `json:"permissions,omitempty"`     // 权限列表
	IsAdmin           bool     `json:"is_admin"`                  // 是否管理员
	LastLoginTime     string   `json:"last_login_time,omitempty"` // 上次登录时间
	LastLoginIP       string   `json:"last_login_ip,omitempty"`   // 上次登录IP
}
