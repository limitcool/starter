package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/pkg/response"
	"github.com/spf13/viper"
)

// 获取系统设置
func GetSystemSettings(c *gin.Context) {
	// 返回当前权限系统设置
	settings := map[string]interface{}{
		"permission": map[string]interface{}{
			"enabled":       global.Config.Permission.Enabled,
			"default_allow": global.Config.Permission.DefaultAllow,
		},
	}

	response.Success(c, settings)
}

// 更新权限系统设置
func UpdatePermissionSettings(c *gin.Context) {
	var req struct {
		Enabled      bool `json:"enabled"`
		DefaultAllow bool `json:"default_allow"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 更新内存中的配置
	global.Config.Permission.Enabled = req.Enabled
	global.Config.Permission.DefaultAllow = req.DefaultAllow

	// 更新配置文件
	v := viper.New()
	v.SetConfigFile(filepath.Join("configs", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		response.InternalServerError(c, "读取配置文件失败: "+err.Error())
		return
	}

	v.Set("permission.enabled", req.Enabled)
	v.Set("permission.default_allow", req.DefaultAllow)

	if err := v.WriteConfig(); err != nil {
		response.InternalServerError(c, "保存配置文件失败: "+err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"message":       "权限系统设置已更新",
		"enabled":       req.Enabled,
		"default_allow": req.DefaultAllow,
	})
}
