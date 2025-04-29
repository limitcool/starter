package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/logger"
)

// AuthService 认证服务
type AuthService struct {
	config *configs.Config
}

// NewAuthService 创建认证服务
func NewAuthService(config *configs.Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

// Claims JWT声明
type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	IsAdmin  bool     `json:"is_admin"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateTokens 生成令牌
func (s *AuthService) GenerateTokens(userID uint, username string, isAdmin bool, roles []string) (*v1.LoginResponse, error) {
	return s.GenerateTokensWithContext(context.Background(), userID, username, isAdmin, roles)
}

// GenerateTokensWithContext 使用上下文生成令牌
func (s *AuthService) GenerateTokensWithContext(ctx context.Context, userID uint, username string, isAdmin bool, roles []string) (*v1.LoginResponse, error) {
	// 获取当前时间
	now := time.Now()

	// 创建访问令牌声明
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(3600) * time.Second)), // 1小时
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "starter-lite",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	}

	// 创建刷新令牌声明
	refreshClaims := &Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(86400) * time.Second)), // 24小时
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "starter-lite",
			Subject:   fmt.Sprintf("%d", userID),
			ID:        fmt.Sprintf("%d", now.UnixNano()),
		},
	}

	// 创建访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte("starter-lite-secret-key"))
	if err != nil {
		logger.ErrorContext(ctx, "生成访问令牌失败", "error", err)
		return nil, errorx.WrapError(err, "生成访问令牌失败")
	}

	// 创建刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("starter-lite-secret-key"))
	if err != nil {
		logger.ErrorContext(ctx, "生成刷新令牌失败", "error", err)
		return nil, errorx.WrapError(err, "生成刷新令牌失败")
	}

	// 返回令牌响应
	return &v1.LoginResponse{
		AccessToken:       accessTokenString,
		RefreshToken:      refreshTokenString,
		ExpiresIn:         3600,
		ExpireTime:        now.Add(time.Duration(3600) * time.Second).Unix(),
		RefreshExpiresIn:  86400,
		RefreshExpireTime: now.Add(time.Duration(86400) * time.Second).Unix(),
		TokenType:         "Bearer",
		UserID:            int64(userID),
		Username:          username,
		Roles:             roles,
		IsAdmin:           isAdmin,
	}, nil
}

// ParseToken 解析令牌
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	return s.ParseTokenWithContext(context.Background(), tokenString)
}

// ParseTokenWithContext 使用上下文解析令牌
func (s *AuthService) ParseTokenWithContext(ctx context.Context, tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("starter-lite-secret-key"), nil
	})

	if err != nil {
		logger.WarnContext(ctx, "解析令牌失败", "error", err)
		return nil, errorx.WrapError(err, "解析令牌失败")
	}

	// 验证令牌
	if !token.Valid {
		logger.WarnContext(ctx, "无效的令牌")
		return nil, errorx.Errorf(errorx.ErrUserTokenError, "无效的令牌")
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok {
		logger.WarnContext(ctx, "无效的令牌声明")
		return nil, errorx.Errorf(errorx.ErrUserTokenError, "无效的令牌声明")
	}

	return claims, nil
}
