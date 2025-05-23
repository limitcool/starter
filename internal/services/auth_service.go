package services

import (
	"context"
	"fmt"
	"time"

	"github.com/limitcool/starter/configs"
	v1 "github.com/limitcool/starter/internal/api/v1"
	"github.com/limitcool/starter/internal/pkg/crypto"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	jwtpkg "github.com/limitcool/starter/internal/pkg/jwt"
)

// AuthService 认证服务，处理通用的认证逻辑
type AuthService struct {
	config *configs.Config
}

// NewAuthService 创建认证服务
func NewAuthService(config *configs.Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

// VerifyPassword 验证用户密码
// 直接使用 crypto.CheckPassword，避免在多个服务中重复实现
func VerifyPassword(password, hashedPassword string) bool {
	return crypto.CheckPassword(hashedPassword, password)
}

// GenerateTokens 生成访问令牌和刷新令牌
// 提取通用的令牌生成逻辑
func (s *AuthService) GenerateTokens(userID uint, username string, userType enum.UserType, roles []string) (*v1.LoginResponse, error) {
	return s.GenerateTokensWithContext(context.Background(), userID, username, userType, roles)
}

// GenerateTokensWithContext 使用上下文生成访问令牌和刷新令牌
func (s *AuthService) GenerateTokensWithContext(ctx context.Context, userID uint, username string, userType enum.UserType, roles []string) (*v1.LoginResponse, error) {
	// 获取配置
	cfg := s.config

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    int64(userID),
		Username:  username,
		UserType:  userType,                      // 用户类型
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
		Roles:     roles,                         // 角色编码
	}

	// 生成刷新令牌
	refreshClaims := &jwtpkg.CustomClaims{
		UserID:    int64(userID),
		Username:  username,
		UserType:  userType,                       // 用户类型
		TokenType: enum.TokenTypeRefresh.String(), // 刷新令牌
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return nil, errorx.WrapError(err, "生成访问令牌失败")
	}

	refreshToken, err := jwtpkg.GenerateTokenWithCustomClaims(refreshClaims, cfg.JwtAuth.RefreshSecret, time.Duration(cfg.JwtAuth.RefreshExpire)*time.Second)
	if err != nil {
		return nil, errorx.WrapError(err, "生成刷新令牌失败")
	}

	return &v1.LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpiresIn:         cfg.JwtAuth.AccessExpire,
		ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
		RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
		RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}, nil
}

// GenerateAccessToken 只生成访问令牌
// 用于刷新令牌时只需要生成新的访问令牌
func (s *AuthService) GenerateAccessToken(userID uint, username string, userType enum.UserType, roles []string) (string, error) {
	return s.GenerateAccessTokenWithContext(context.Background(), userID, username, userType, roles)
}

// GenerateAccessTokenWithContext 使用上下文只生成访问令牌
func (s *AuthService) GenerateAccessTokenWithContext(ctx context.Context, userID uint, username string, userType enum.UserType, roles []string) (string, error) {
	// 获取配置
	cfg := s.config

	// 生成访问令牌
	accessClaims := &jwtpkg.CustomClaims{
		UserID:    int64(userID),
		Username:  username,
		UserType:  userType,                      // 用户类型
		TokenType: enum.TokenTypeAccess.String(), // 访问令牌
		Roles:     roles,                         // 角色编码
	}

	accessToken, err := jwtpkg.GenerateTokenWithCustomClaims(accessClaims, cfg.JwtAuth.AccessSecret, time.Duration(cfg.JwtAuth.AccessExpire)*time.Second)
	if err != nil {
		return "", errorx.WrapError(err, fmt.Sprintf("生成用户 %s 的访问令牌失败", username))
	}

	return accessToken, nil
}

// ParseRefreshToken 解析刷新令牌
func (s *AuthService) ParseRefreshToken(refreshToken string) (*jwtpkg.CustomClaims, error) {
	return s.ValidateRefreshTokenWithContext(context.Background(), refreshToken)
}

// ValidateRefreshToken 验证刷新令牌
func (s *AuthService) ValidateRefreshToken(refreshToken string) (*jwtpkg.CustomClaims, error) {
	return s.ValidateRefreshTokenWithContext(context.Background(), refreshToken)
}

// ValidateRefreshTokenWithContext 使用上下文验证刷新令牌
func (s *AuthService) ValidateRefreshTokenWithContext(ctx context.Context, refreshToken string) (*jwtpkg.CustomClaims, error) {
	// 获取配置
	cfg := s.config

	// 验证刷新令牌
	claims, err := jwtpkg.ParseTokenWithCustomClaims(refreshToken, cfg.JwtAuth.RefreshSecret)
	if err != nil {
		return nil, errorx.WrapError(err, "验证刷新令牌失败")
	}

	// 检查令牌类型
	if claims.TokenType != enum.TokenTypeRefresh.String() {
		tokenErr := errorx.Errorf(errorx.ErrUserTokenError, "无效的令牌类型")
		return nil, errorx.WrapError(tokenErr, "")
	}

	return claims, nil
}

// CreateLoginResponse 创建登录响应
func (s *AuthService) CreateLoginResponse(accessToken, refreshToken string) *v1.LoginResponse {
	return s.CreateLoginResponseWithContext(context.Background(), accessToken, refreshToken)
}

// CreateLoginResponseWithContext 使用上下文创建登录响应
func (s *AuthService) CreateLoginResponseWithContext(ctx context.Context, accessToken, refreshToken string) *v1.LoginResponse {
	// 获取配置
	cfg := s.config

	return &v1.LoginResponse{
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpiresIn:         cfg.JwtAuth.AccessExpire,
		ExpireTime:        time.Now().Add(time.Duration(cfg.JwtAuth.AccessExpire) * time.Second).Unix(),
		RefreshExpiresIn:  cfg.JwtAuth.RefreshExpire,
		RefreshExpireTime: time.Now().Add(time.Duration(cfg.JwtAuth.RefreshExpire) * time.Second).Unix(),
	}
}
