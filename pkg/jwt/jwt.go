package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenExpireDuration 过期时间
const (
	TokenExpireDuration = time.Hour * 2
)

// GenerateToken 生成JWT Token
func GenerateToken(claims jwt.MapClaims, secret string, expireDuration time.Duration) (string, error) {
	// 设置过期时间
	if expireDuration <= 0 {
		expireDuration = TokenExpireDuration
	}
	claims["exp"] = time.Now().Add(expireDuration).Unix()
	claims["iat"] = time.Now().Unix()

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	return token.SignedString([]byte(secret))
}

// ParseToken 解析和校验token
func ParseToken(token, secret string) (*jwt.MapClaims, error) {
	// 解析token
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if jwtToken != nil {
		// 校验token
		if Claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
			return &Claims, nil
		}
	}
	return nil, err
}
