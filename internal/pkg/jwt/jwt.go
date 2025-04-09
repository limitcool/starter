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

// GenerateTokenWithCustomClaims 使用CustomClaims结构体生成JWT Token
// 相比GenerateToken，此函数支持直接使用封装好的CustomClaims结构体
// 提供更好的类型安全和代码可读性
func GenerateTokenWithCustomClaims(claims *CustomClaims, secret string, expireDuration time.Duration) (string, error) {
	// 设置过期时间
	if expireDuration <= 0 {
		expireDuration = TokenExpireDuration
	}

	// 设置过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expireDuration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

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

// ParseTokenWithCustomClaims 解析和校验token，返回CustomClaims结构体
// 相比ParseToken，此函数直接返回CustomClaims结构体
// 方便后续使用结构体字段而不是通过map访问
func ParseTokenWithCustomClaims(token, secret string) (*CustomClaims, error) {
	// 先解析为MapClaims
	mapClaims, err := ParseToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 转换为CustomClaims
	return FromMapClaims(*mapClaims), nil
}
