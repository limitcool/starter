package errorx

import (
	"net/http"

	"github.com/epkgs/i18n"
)

func init() {
	fileI18n.LoadTranslations()
}

var fileI18n = i18n.NewCatalog("file")

var (
	ErrFileNotFound            = defineErr[none](fileI18n, 4000, "file does not exist", http.StatusNotFound)                     // 文件不存在
	ErrFileUpload              = defineErr[none](fileI18n, 4001, "file upload failed", http.StatusInternalServerError)           // 文件上传失败
	ErrFileDelete              = defineErr[none](fileI18n, 4002, "file deletion failed", http.StatusInternalServerError)         // 文件删除失败
	ErrFileUpdate              = defineErr[none](fileI18n, 4003, "file update failed", http.StatusInternalServerError)           // 文件更新失败
	ErrFileDownload            = defineErr[none](fileI18n, 4004, "file download failed", http.StatusInternalServerError)         // 文件下载失败
	ErrGetUploadURL            = defineErr[none](fileI18n, 4005, "get upload url failed", http.StatusInternalServerError)        // 获取上传URL失败
	ErrFileCreate              = defineErr[none](fileI18n, 4006, "file create failed", http.StatusInternalServerError)           // 创建文件记录失败
	ErrFileNotExist            = defineErr[none](fileI18n, 4007, "file not exist", http.StatusNotFound)                          // 文件记录不存在
	ErrFileVerify              = defineErr[none](fileI18n, 4008, "file verify failed", http.StatusInternalServerError)           // 验证文件失败
	ErrFileUploadNotComplete   = defineErr[none](fileI18n, 4009, "file upload not complete", http.StatusNotFound)                // 文件上传未完成
	ErrFileGenerateDownloadURL = defineErr[none](fileI18n, 4010, "generate download url failed", http.StatusInternalServerError) // 生成下载URL失败
	ErrFileUpdateRecord        = defineErr[none](fileI18n, 4011, "file update record failed", http.StatusInternalServerError)    // 更新文件记录失败
	ErrFileIDEmpty             = defineErr[none](fileI18n, 4012, "file id can not be empty", http.StatusBadRequest)              // 文件ID不能为空
	ErrGetUploadFile           = defineErr[none](fileI18n, 4013, "get upload file failed", http.StatusBadRequest)                // 获取上传文件失败
	ErrOpenUploadFile          = defineErr[none](fileI18n, 4014, "open upload file failed", http.StatusBadRequest)               // 打开上传文件失败
)
