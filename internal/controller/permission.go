package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/pkg/response"
)

// GetPermissions 获取权限列表
func GetPermissions(c *gin.Context) {
	var permissions []model.Permission
	if err := global.DB.Find(&permissions).Error; err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, permissions)
}

// GetPermission 获取权限详情
func GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var permission model.Permission
	if err := global.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := global.DB.Create(&permission).Error; err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	permission.ID = uint(id)
	if err := global.DB.Model(&model.Permission{}).Where("id = ?", id).Updates(permission).Error; err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的权限ID")
		return
	}

	// 删除权限
	if err := global.DB.Delete(&model.Permission{}, id).Error; err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
