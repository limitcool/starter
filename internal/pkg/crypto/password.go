package crypto

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 将密码使用bcrypt算法加密
func HashPassword(password string) (string, error) {
	return HashPasswordWithContext(context.Background(), password)
}

// HashPasswordWithContext 使用上下文将密码使用bcrypt算法加密
func HashPasswordWithContext(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码是否匹配
func CheckPassword(hashedPassword, password string) bool {
	return CheckPasswordWithContext(context.Background(), hashedPassword, password)
}

// CheckPasswordWithContext 使用上下文验证密码是否匹配
func CheckPasswordWithContext(ctx context.Context, hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
