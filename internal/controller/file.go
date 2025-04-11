package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/enum"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/cast"
)

// FileController 文件控制器
type FileController struct {
	fileService *services.FileService
}

// NewFileController 创建文件控制器
func NewFileController(storage *storage.Storage) *FileController {
	return &FileController{
		fileService: services.NewFileService(storage),
	}
}

// UploadFile 上传文件
func (ctrl *FileController) UploadFile(c *gin.Context) {
	// 获取当前用户ID
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取用户类型
	userType, ok := c.Get("user_type")
	if !ok {
		userType = enum.UserTypeSysUser // 默认为系统用户
	}

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取文件类型参数
	fileType := c.DefaultPostForm("type", model.FileTypeOther)

	// 上传文件，传入用户类型
	fileModel, err := ctrl.fileService.UploadFile(c, file, fileType, cast.ToInt64(userID), enum.UserType(cast.ToInt8(userType)))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, fileModel)
}

// GetFile 获取文件信息
func (ctrl *FileController) GetFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取当前用户ID和类型
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	userType, ok := c.Get("user_type")
	if !ok {
		userType = enum.UserTypeSysUser // 默认为系统用户
	}

	file, err := ctrl.fileService.GetFile(id, cast.ToInt64(userID), enum.UserType(cast.ToUint8(userType)))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, file)
}

// DeleteFile 删除文件
func (ctrl *FileController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取当前用户ID和类型
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	userType, ok := c.Get("user_type")
	if !ok {
		userType = "sys_user" // 默认为系统用户
	}

	err := ctrl.fileService.DeleteFile(id, cast.ToInt64(userID), enum.UserType(cast.ToUint8(userType)))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success[interface{}](c, nil)
}

// DownloadFile 下载文件
func (ctrl *FileController) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取当前用户ID和类型
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	userType, ok := c.Get("user_type")
	if !ok {
		userType = enum.UserTypeSysUser // 默认为系统用户
	}

	// 获取文件内容
	fileStream, contentType, err := ctrl.fileService.GetFileContent(id, cast.ToInt64(userID), enum.UserType(cast.ToUint8(userType)))
	if err != nil {
		response.Error(c, err)
		return
	}
	defer fileStream.Close()

	// 获取文件信息以设置文件名
	file, err := ctrl.fileService.GetFile(id, cast.ToInt64(userID), enum.UserType(cast.ToUint8(userType)))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 设置响应头
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))

	// 将文件流写入响应
	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		response.Error(c, errorx.ErrFileStroage)
		return
	}
}

// UpdateUserAvatar 更新用户头像
func (ctrl *FileController) UpdateUserAvatar(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取头像文件
	avatar, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 更新头像
	avatarURL, err := ctrl.fileService.UpdateUserAvatar(cast.ToInt64(userID), avatar)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, avatarURL)
}

// UpdateSysUserAvatar 更新系统用户头像
func (ctrl *FileController) UpdateSysUserAvatar(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		response.Error(c, errorx.ErrUserNoLogin)
		return
	}

	// 获取头像文件
	avatar, err := c.FormFile("avatar")
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 更新头像
	avatarURL, err := ctrl.fileService.UpdateSysUserAvatar(cast.ToInt64(userID), avatar)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, avatarURL)
}
