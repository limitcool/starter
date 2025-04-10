package controller

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/core"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/pkg/errorx"
	"github.com/limitcool/starter/internal/storage/sqldb"
	"github.com/spf13/viper"
)

func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

type PermissionController struct {
}

// 更新权限系统设置
func (pc *PermissionController) UpdatePermissionSettings(c *gin.Context) {
	var req struct {
		Enabled      bool `json:"enabled"`
		DefaultAllow bool `json:"default_allow"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 获取配置
	config := core.Instance().Config()

	// 更新内存中的配置
	config.Permission.Enabled = req.Enabled
	config.Permission.DefaultAllow = req.DefaultAllow

	// 更新配置文件
	v := viper.New()
	v.SetConfigFile(filepath.Join("configs", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		response.Error(c, err)
		return
	}

	v.Set("permission.enabled", req.Enabled)
	v.Set("permission.default_allow", req.DefaultAllow)

	if err := v.WriteConfig(); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, map[string]interface{}{
		"message":       "权限系统设置已更新",
		"enabled":       req.Enabled,
		"default_allow": req.DefaultAllow,
	})
}

// GetPermissions 获取权限列表
func (pc *PermissionController) GetPermissions(c *gin.Context) {
	var permissions []model.Permission
	db := sqldb.GetDB()
	if err := db.Find(&permissions).Error; err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, permissions)
}

// GetPermission 获取权限详情
func (pc *PermissionController) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	var permission model.Permission
	db := sqldb.GetDB()
	if err := db.Where("id = ?", id).First(&permission).Error; err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, err)
		return
	}

	db := sqldb.GetDB()
	if err := db.Create(&permission).Error; err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
func (pc *PermissionController) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, err)
		return
	}

	permission.ID = uint(id)
	db := sqldb.GetDB()
	if err := db.Model(&model.Permission{}).Where("id = ?", id).Updates(permission).Error; err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}

// DeletePermission 删除权限
func (pc *PermissionController) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errorx.ErrInvalidParams)
		return
	}

	// 删除权限
	db := sqldb.GetDB()
	if err := db.Delete(&model.Permission{}, id).Error; err != nil {
		response.Error(c, err)
		return
	}

	response.Success[any](c, nil)
}
