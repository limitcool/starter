package errorx

import "net/http"

// 文件错误码 (4000-4999)
const (
	ErrorFileNotFoundCode = 4000 // 文件不存在
	ErrorFileUploadCode   = 4001 // 文件上传失败
	ErrorFileDeleteCode   = 4002 // 文件删除失败
	ErrorFileUpdateCode   = 4003 // 文件更新失败
	ErrorFileDownloadCode = 4004 // 文件下载失败
)

// 预定义文件错误实例
var (
	ErrFileNotFound = NewAppError(ErrorFileNotFoundCode, "文件不存在", http.StatusNotFound)
	ErrFileUpload   = NewAppError(ErrorFileUploadCode, "文件上传失败", http.StatusInternalServerError)
	ErrFileDelete   = NewAppError(ErrorFileDeleteCode, "文件删除失败", http.StatusInternalServerError)
	ErrFileUpdate   = NewAppError(ErrorFileUpdateCode, "文件更新失败", http.StatusInternalServerError)
	ErrFileDownload = NewAppError(ErrorFileDownloadCode, "文件下载失败", http.StatusInternalServerError)
)
