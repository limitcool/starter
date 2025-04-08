package code

var MsgFlags = map[int]string{
	Success:           "成功",
	InvalidParams:     "请求参数错误",
	ErrorUnknown:      "服务器开小差啦，稍后再来试一试",
	ErrorNotExistCert: "不存在的认证类型",
	ErrorNotFound:     "资源不存在",
	ErrorDatabase:     "数据库操作失败",
	ErrorInternal:     "服务器内部错误",
	Error:             "错误",
	ErrorParam:        "参数错误",

	DatabaseInsertError: "数据库插入失败",
	DatabaseDeleteError: "数据库删除失败",
	DatabaseQueryError:  "数据库查询失败",

	UserNoLogin:             "用户未登录",
	UserNotFound:            "用户不存在",
	UserPasswordError:       "密码错误",
	UserNotVerify:           "用户未验证",
	UserLocked:              "用户已锁定",
	UserDisabled:            "用户已禁用",
	UserExpired:             "用户已过期",
	UserAlreadyExists:       "用户已存在",
	UserNameOrPasswordError: "用户名或密码错误",
	UserAuthFailed:          "认证失败",
	UserNoPermission:        "没有权限",
	UserPasswordErr:         "密码错误",
	UserNotExist:            "用户不存在",
	UserTokenError:          "登录凭证无效",
	UserTokenExpired:        "登录已过期，请重新登录",

	AccessDenied: "访问被拒绝",

	ErrorCache:        "缓存操作失败",
	ErrorCacheTimeout: "缓存操作超时",
	ErrorCacheKey:     "缓存键不存在",
	ErrorCacheValue:   "缓存值错误",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ErrorUnknown]
}
