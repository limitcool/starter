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
