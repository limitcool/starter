package jwtx

import (
	"testing"

	"github.com/spf13/cast"
)

func TestParseToken(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name    string
		token   string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   "",
			secret:  "",
			wantErr: false,
		},
	}

	// 测试
	for _, tc := range testCases {
		claims, err := ParseToken(tc.token, tc.secret)
		for k, v := range *claims {
			t.Log(k, v)
		}
		t.Error(err)
		t.Log(cast.ToInt((*claims)["user_id"]))

	}
}
