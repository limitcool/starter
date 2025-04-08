package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name           string
		claims         jwt.MapClaims
		secret         string
		expireDuration time.Duration
		wantErr        bool
	}{
		{
			name: "Valid token generation",
			claims: jwt.MapClaims{
				"user_id":  1,
				"username": "testuser",
			},
			secret:         "testsecret",
			expireDuration: time.Hour,
			wantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := GenerateToken(tc.claims, tc.secret, tc.expireDuration)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 测试解析生成的token
				claims, err := ParseToken(token, tc.secret)
				assert.NoError(t, err)
				assert.Equal(t, tc.claims["user_id"], (*claims)["user_id"])
				assert.Equal(t, tc.claims["username"], (*claims)["username"])
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	// 生成测试token
	claims := jwt.MapClaims{
		"user_id":  123,
		"username": "testuser",
	}
	secret := "testsecret"
	token, err := GenerateToken(claims, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate token for test: %v", err)
	}

	// 测试用例
	testCases := []struct {
		name    string
		token   string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   token,
			secret:  secret,
			wantErr: false,
		},
		{
			name:    "Invalid secret",
			token:   token,
			secret:  "wrongsecret",
			wantErr: true,
		},
		{
			name:    "Invalid token",
			token:   "invalidtoken",
			secret:  secret,
			wantErr: true,
		},
	}

	// 测试
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			claims, err := ParseToken(tc.token, tc.secret)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				if claims != nil {
					t.Log("user_id:", cast.ToInt((*claims)["user_id"]))
					assert.Equal(t, 123, cast.ToInt((*claims)["user_id"]))
				}
			}
		})
	}
}
