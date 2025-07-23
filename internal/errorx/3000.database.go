package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

func init() {
	dbI18n.LoadTranslations()
}

var dbI18n = i18n.NewCatalog("database")

var (
	ErrDatabase            = defineErr[none](dbI18n, 3000, "database error", http.StatusInternalServerError)             // 数据库错误
	ErrDatabaseQuery       = defineErr[none](dbI18n, 3001, "database query error", http.StatusInternalServerError)       // 数据库查询错误
	ErrDatabaseInsert      = defineErr[none](dbI18n, 3002, "database insert error", http.StatusInternalServerError)      // 数据库插入错误
	ErrDatabaseUpdate      = defineErr[none](dbI18n, 3003, "database update error", http.StatusInternalServerError)      // 数据库更新错误
	ErrDatabaseDelete      = defineErr[none](dbI18n, 3004, "database delete error", http.StatusInternalServerError)      // 数据库删除错误
	ErrDatabaseConnection  = defineErr[none](dbI18n, 3005, "database connection error", http.StatusInternalServerError)  // 数据库连接错误
	ErrDatabaseTransaction = defineErr[none](dbI18n, 3006, "database transaction error", http.StatusInternalServerError) // 数据库事务错误
	ErrQueryParamEmpty     = defineErr[none](dbI18n, 3007, "query parameter cannot be empty", http.StatusBadRequest)     // 查询参数不能为空
	ErrRecordNotExist      = defineErr[none](dbI18n, 3008, "record does not exist", http.StatusNotFound)                 // 记录不存在
	ErrQueryUser           = defineErr[none](dbI18n, 3009, "query user failed", http.StatusInternalServerError)          // 查询用户失败
	ErrQueryUserAvatar     = defineErr[none](dbI18n, 3010, "query user avatar failed", http.StatusInternalServerError)   // 查询用户头像失败
	ErrCheckUserExist      = defineErr[none](dbI18n, 3011, "check user exist failed", http.StatusInternalServerError)    // 检查用户是否存在失败
	ErrQueryUserList       = defineErr[none](dbI18n, 3012, "query user list failed", http.StatusInternalServerError)     // 查询用户列表失败
	ErrQueryUserTotal      = defineErr[none](dbI18n, 3013, "query user total failed", http.StatusInternalServerError)    // 查询用户总数失败
	ErrQueryFile           = defineErr[none](dbI18n, 3014, "query file failed", http.StatusBadRequest)                   // 查询文件失败
	ErrQueryUserFileList   = defineErr[none](dbI18n, 3015, "query user file list failed", http.StatusBadRequest)         // 查询用户文件列表失败
	ErrQueryUserFileTotal  = defineErr[none](dbI18n, 3016, "query user file total failed", http.StatusBadRequest)        // 查询用户文件总数失败
	ErrQueryFileList       = defineErr[none](dbI18n, 3017, "query file list failed", http.StatusBadRequest)              // 查询文件列表失败
	ErrQueryFileTotal      = defineErr[none](dbI18n, 3018, "query file total failed", http.StatusBadRequest)             // 查询文件总数失败
)
