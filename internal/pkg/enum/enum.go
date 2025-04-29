package enum

// TokenType 令牌类型
type TokenType uint8

const (
	TokenTypeAccess  TokenType = iota + 1 // 访问令牌
	TokenTypeRefresh                      // 刷新令牌
)

func (t TokenType) String() string {
	return []string{"access_token", "refresh_token"}[t-1]
}
