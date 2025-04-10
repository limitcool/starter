package v1

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求参数
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
