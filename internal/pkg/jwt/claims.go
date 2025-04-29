package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims 自定义JWT Claims结构体
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	IsAdmin   bool     `json:"is_admin"`             // 是否是管理员
	TokenType string   `json:"token_type,omitempty"` // access_token 或 refresh_token
	RoleIDs   []uint   `json:"role_ids,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

// ToMapClaims 将CustomClaims转换为jwt.MapClaims
func (c *CustomClaims) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"user_id":    c.UserID,
		"username":   c.Username,
		"is_admin":   c.IsAdmin,
		"token_type": c.TokenType,
		"roles":      c.Roles,
		"role_ids":   c.RoleIDs,
		"exp":        c.ExpiresAt.Unix(),
		"iat":        c.IssuedAt.Unix(),
	}
}

// FromMapClaims 从jwt.MapClaims转换为CustomClaims
func FromMapClaims(claims jwt.MapClaims) *CustomClaims {
	customClaims := &CustomClaims{}

	// 用户ID
	if userID, ok := claims["user_id"].(float64); ok {
		customClaims.UserID = int64(userID)
	}

	// 用户名
	if username, ok := claims["username"].(string); ok {
		customClaims.Username = username
	}

	// 是否管理员
	if isAdmin, ok := claims["is_admin"].(bool); ok {
		customClaims.IsAdmin = isAdmin
	}

	// 令牌类型
	if tokenType, ok := claims["token_type"].(string); ok {
		customClaims.TokenType = tokenType
	} else if tokenType, ok := claims["type"].(string); ok { // 兼容旧版
		customClaims.TokenType = tokenType
	}

	// 角色代码
	if rolesInterface, ok := claims["roles"].([]any); ok {
		roles := make([]string, len(rolesInterface))
		for i, v := range rolesInterface {
			if role, ok := v.(string); ok {
				roles[i] = role
			}
		}
		customClaims.Roles = roles
	}

	// 角色ID
	if roleIDsInterface, ok := claims["role_ids"].([]any); ok {
		roleIDs := make([]uint, len(roleIDsInterface))
		for i, v := range roleIDsInterface {
			if roleID, ok := v.(float64); ok {
				roleIDs[i] = uint(roleID)
			}
		}
		customClaims.RoleIDs = roleIDs
	}

	return customClaims
}
