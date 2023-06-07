package jwtx

import "github.com/golang-jwt/jwt/v4"

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
