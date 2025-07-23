package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

func init() {
	userI18n.LoadTranslations()
}

var userI18n = i18n.NewCatalog("user")

var (
	ErrUserNotFound            = defineErr[none](userI18n, 2000, "user not found", http.StatusNotFound)                                         // 用户不存在
	ErrInvalidCredentials      = defineErr[none](userI18n, 2001, "invalid credentials", http.StatusUnauthorized)                                // 无效的凭证
	ErrUserDisabled            = defineErr[struct{ Name string }](userI18n, 2002, "user {{.Name}} is disabled", http.StatusForbidden)           // 用户 {{.Name}} 已被禁用
	ErrUserExists              = defineErr[struct{ Name string }](userI18n, 2003, "user {{.Name}} already exists", http.StatusConflict)         // 用户 {{.Name}} 已存在
	ErrPasswordExpired         = defineErr[none](userI18n, 2004, "password expired", http.StatusForbidden)                                      // 密码已过期
	ErrUserPassword            = defineErr[struct{ Name string }](userI18n, 2005, "user {{.Name}} password incorrect", http.StatusUnauthorized) // 用户 {{.Name}} 的密码错误
	ErrUserTokenError          = defineErr[none](userI18n, 2006, "wrong user token", http.StatusUnauthorized)                                   // 用户令牌错误
	ErrUserNotLogin            = defineErr[none](userI18n, 2007, "user is not logged in", http.StatusUnauthorized)                              // 用户未登录
	ErrGenVisitToken           = defineErr[none](userI18n, 2008, "generate visit token failed", http.StatusInternalServerError)                 // 生成访问令牌失败
	ErrGenRefreshToken         = defineErr[none](userI18n, 2009, "generate refresh token failed", http.StatusInternalServerError)               // 生成刷新令牌失败
	ErrParseToken              = defineErr[none](userI18n, 2010, "parse token failed", http.StatusUnauthorized)                                 // 解析令牌失败
	ErrInvalidToken            = defineErr[none](userI18n, 2011, "invalid token", http.StatusUnauthorized)                                      // 无效的令牌
	ErrInvalidTokenClaim       = defineErr[none](userI18n, 2012, "invalid token claim", http.StatusUnauthorized)                                // 无效的令牌声明
	ErrPasswordEncrypt         = defineErr[none](userI18n, 2013, "password encrypt failed", http.StatusInternalServerError)                     // 密码加密失败
	ErrPasswordDecrypt         = defineErr[none](userI18n, 2014, "password decrypt failed", http.StatusInternalServerError)                     // 密码解密失败
	ErrOldPasswordError        = defineErr[none](userI18n, 2015, "old password error", http.StatusUnauthorized)                                 // 旧密码错误
	ErrUserNameOrPasswordEmpty = defineErr[none](userI18n, 2016, "username or password empty", http.StatusBadRequest)                           // 用户名或密码不能为空
	ErrPassword                = defineErr[none](userI18n, 2017, "password error", http.StatusUnauthorized)                                     // 密码错误
)
