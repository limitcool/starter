package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/model"
	"github.com/limitcool/starter/internal/services"
)

// GetPermissions 获取权限列表
func GetPermissions(c *gin.Context) {
	var permissions []model.Permission
	db := services.Instance().GetDB()
	if err := db.Find(&permissions).Error; err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, permissions)
}

// GetPermission 获取权限详情
func GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的权限ID")
		return
	}

	var permission model.Permission
	db := services.Instance().GetDB()
	if err := db.Where("id = ?", id).First(&permission).Error; err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, permission)
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	db := services.Instance().GetDB()
	if err := db.Create(&permission).Error; err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的权限ID")
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	permission.ID = uint(id)
	db := services.Instance().GetDB()
	if err := db.Model(&model.Permission{}).Where("id = ?", id).Updates(permission).Error; err != nil {
		response.ServerError(c)
		return
	}

	response.Success[any](c, nil)
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ParamError(c, "无效的权限ID")
		return
	}

	// 删除权限
	db := services.Instance().GetDB()
	if err := db.Delete(&model.Permission{}, id).Error; err != nil {
		response.ServerError(c)
		return
	}

	response.Success[any](c, nil)
}
