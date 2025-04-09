package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
	"github.com/limitcool/starter/pkg/apiresponse"
	"github.com/limitcool/starter/pkg/errors"
	"github.com/limitcool/starter/pkg/storage"
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
// @Summary 上传文件
// @Description 上传文件到服务器
// @Tags 文件
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type formData string false "文件类型(image/avatar/document/video/audio/other)" default(other)
// @Success 200 {object} apiresponse.Response{data=model.File}
// @Router /api/files/upload [post]
func (ctrl *FileController) UploadFile(c *gin.Context) {
	// 获取当前用户ID
	var userID string
	user, exists := c.Get("user")
	if exists {
		if u, ok := user.(*model.User); ok {
			userID = fmt.Sprintf("%d", u.ID)
		} else if su, ok := user.(*model.SysUser); ok {
			userID = fmt.Sprintf("%d", su.ID)
		}
	}

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		apiresponse.ParamError(c, "获取上传文件失败")
		return
	}

	// 获取文件类型参数
	fileType := c.DefaultPostForm("type", model.FileTypeOther)

	// 上传文件
	fileModel, err := ctrl.fileService.UploadFile(c, file, fileType, userID)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[*model.File](c, fileModel)
}

// GetFile 获取文件信息
// @Summary 获取文件信息
// @Description 获取文件详细信息
// @Tags 文件
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} apiresponse.Response{data=model.File}
// @Router /api/files/{id} [get]
func (ctrl *FileController) GetFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apiresponse.ParamError(c, "文件ID不能为空")
		return
	}

	file, err := ctrl.fileService.GetFile(id)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[*model.File](c, file)
}

// DeleteFile 删除文件
// @Summary 删除文件
// @Description 删除指定文件
// @Tags 文件
// @Produce json
// @Param id path string true "文件ID"
// @Success 200 {object} apiresponse.Response
// @Router /api/files/{id} [delete]
func (ctrl *FileController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apiresponse.ParamError(c, "文件ID不能为空")
		return
	}

	err := ctrl.fileService.DeleteFile(id)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[interface{}](c, nil)
}

// DownloadFile 下载文件
// @Summary 下载文件
// @Description 下载指定文件
// @Tags 文件
// @Produce octet-stream
// @Param id path string true "文件ID"
// @Router /api/files/{id}/download [get]
func (ctrl *FileController) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		apiresponse.ParamError(c, "文件ID不能为空")
		return
	}

	// 获取文件内容
	fileStream, contentType, err := ctrl.fileService.GetFileContent(id)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}
	defer fileStream.Close()

	// 获取文件信息以设置文件名
	file, err := ctrl.fileService.GetFile(id)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	// 设置响应头
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))

	// 将文件流写入响应
	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		apiresponse.HandleError(c, errors.NewStorageError("下载", file.Path, err))
		return
	}
}

// UpdateUserAvatar 更新用户头像
// @Summary 更新用户头像
// @Description 更新当前用户的头像
// @Tags 用户
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "头像文件"
// @Success 200 {object} apiresponse.Response{data=string}
// @Router /api/users/avatar [put]
func (ctrl *FileController) UpdateUserAvatar(c *gin.Context) {
	// 获取当前用户
	user, exists := c.Get("user")
	if !exists {
		apiresponse.Unauthorized(c, "未登录")
		return
	}

	// 获取用户ID
	var userID string
	if u, ok := user.(*model.User); ok {
		userID = strconv.FormatUint(uint64(u.ID), 10)
	} else {
		apiresponse.ParamError(c, "用户类型无效")
		return
	}

	// 获取头像文件
	avatar, err := c.FormFile("avatar")
	if err != nil {
		apiresponse.ParamError(c, "获取头像文件失败")
		return
	}

	// 更新头像
	avatarURL, err := ctrl.fileService.UpdateUserAvatar(userID, avatar)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[string](c, avatarURL)
}

// UpdateSysUserAvatar 更新系统用户头像
// @Summary 更新系统用户头像
// @Description 更新当前系统用户的头像
// @Tags 系统用户
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "头像文件"
// @Success 200 {object} apiresponse.Response{data=string}
// @Router /api/sysusers/avatar [put]
func (ctrl *FileController) UpdateSysUserAvatar(c *gin.Context) {
	// 获取当前用户
	user, exists := c.Get("user")
	if !exists {
		apiresponse.Unauthorized(c, "未登录")
		return
	}

	// 获取用户ID
	var userID string
	if su, ok := user.(*model.SysUser); ok {
		userID = strconv.FormatUint(uint64(su.ID), 10)
	} else {
		apiresponse.ParamError(c, "用户类型无效")
		return
	}

	// 获取头像文件
	avatar, err := c.FormFile("avatar")
	if err != nil {
		apiresponse.ParamError(c, "获取头像文件失败")
		return
	}

	// 更新头像
	avatarURL, err := ctrl.fileService.UpdateSysUserAvatar(userID, avatar)
	if err != nil {
		apiresponse.HandleError(c, err)
		return
	}

	apiresponse.Success[string](c, avatarURL)
}
