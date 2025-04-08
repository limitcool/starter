package code

const (
	Success = 0
)

const (
	InvalidParams = 10000 + iota
	ErrorUnknown
	ErrorNotExistCert
	ErrorNotFound // 资源不存在
	ErrorDatabase // 数据库操作失败
)

const (
	DatabaseInsertError = 20000 + iota // 数据库插入错误
	DatabaseDeleteError                // 数据库删除失败
	DatabaseQueryError                 // 数据库查询失败
)

const (
	UserNoLogin             = 30000 + iota // 用户未登录
	UserNotFound                           // 用户不存在
	UserPasswordError                      // 用户密码错误
	UserNotVerify                          // 用户未验证
	UserLocked                             // 用户被锁定
	UserDisabled                           // 用户被关闭
	UserExpired                            // 用户过期
	UserAlreadyExists                      // 用户已存在
	UserNameOrPasswordError                // 用户名或密码错误
	UserAuthFailed                         // 用户鉴权失败
	UserNoPermission                       //用户无权访问
)
