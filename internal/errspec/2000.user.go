package errspec

import (
	"net/http"

	"github.com/epkgs/i18n"
	"github.com/limitcool/starter/internal/pkg/errorx"
)

func init() {
	userI18n.LoadTranslations()
}

var userI18n = i18n.NewCatalog("user")

var (
	ErrUserNotFound            = errorx.DefineSimple(userI18n, 2000, "user not found", http.StatusNotFound)                                         // 用户不存在
	ErrInvalidCredentials      = errorx.DefineSimple(userI18n, 2001, "invalid credentials", http.StatusUnauthorized)                                // 无效的凭证
	ErrUserDisabled            = errorx.Define[struct{ Name string }](userI18n, 2002, "user {{.Name}} is disabled", http.StatusForbidden)           // 用户 {{.Name}} 已被禁用
	ErrUserExists              = errorx.Define[struct{ Name string }](userI18n, 2003, "user {{.Name}} already exists", http.StatusConflict)         // 用户 {{.Name}} 已存在
	ErrPasswordExpired         = errorx.DefineSimple(userI18n, 2004, "password expired", http.StatusForbidden)                                      // 密码已过期
	ErrUserPassword            = errorx.Define[struct{ Name string }](userI18n, 2005, "user {{.Name}} password incorrect", http.StatusUnauthorized) // 用户 {{.Name}} 的密码错误
	ErrUserTokenError          = errorx.DefineSimple(userI18n, 2006, "wrong user token", http.StatusUnauthorized)                                   // 用户令牌错误
	ErrUserNotLogin            = errorx.DefineSimple(userI18n, 2007, "user is not logged in", http.StatusUnauthorized)                              // 用户未登录
	ErrGenVisitToken           = errorx.DefineSimple(userI18n, 2008, "generate visit token failed", http.StatusInternalServerError)                 // 生成访问令牌失败
	ErrGenRefreshToken         = errorx.DefineSimple(userI18n, 2009, "generate refresh token failed", http.StatusInternalServerError)               // 生成刷新令牌失败
	ErrParseToken              = errorx.DefineSimple(userI18n, 2010, "parse token failed", http.StatusUnauthorized)                                 // 解析令牌失败
	ErrInvalidToken            = errorx.DefineSimple(userI18n, 2011, "invalid token", http.StatusUnauthorized)                                      // 无效的令牌
	ErrInvalidTokenClaim       = errorx.DefineSimple(userI18n, 2012, "invalid token claim", http.StatusUnauthorized)                                // 无效的令牌声明
	ErrPasswordEncrypt         = errorx.DefineSimple(userI18n, 2013, "password encrypt failed", http.StatusInternalServerError)                     // 密码加密失败
	ErrPasswordDecrypt         = errorx.DefineSimple(userI18n, 2014, "password decrypt failed", http.StatusInternalServerError)                     // 密码解密失败
	ErrOldPasswordError        = errorx.DefineSimple(userI18n, 2015, "old password error", http.StatusUnauthorized)                                 // 旧密码错误
	ErrUserNameOrPasswordEmpty = errorx.DefineSimple(userI18n, 2016, "username or password empty", http.StatusBadRequest)                           // 用户名或密码不能为空
	ErrPassword                = errorx.DefineSimple(userI18n, 2017, "password error", http.StatusUnauthorized)                                     // 密码错误
)
