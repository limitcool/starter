package controller

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/services"
	"github.com/spf13/viper"
)

// 获取系统设置
func GetSystemSettings(c *gin.Context) {
	// 获取配置
	config := services.Instance().GetConfig()

	// 返回当前权限系统设置
	settings := map[string]interface{}{
		"permission": map[string]interface{}{
			"enabled":       config.Permission.Enabled,
			"default_allow": config.Permission.DefaultAllow,
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
		response.ParamError(c, err.Error())
		return
	}

	// 获取配置
	config := services.Instance().GetConfig()

	// 更新内存中的配置
	config.Permission.Enabled = req.Enabled
	config.Permission.DefaultAllow = req.DefaultAllow

	// 更新配置文件
	v := viper.New()
	v.SetConfigFile(filepath.Join("configs", "config.yaml"))

	if err := v.ReadInConfig(); err != nil {
		response.ServerError(c)
		return
	}

	v.Set("permission.enabled", req.Enabled)
	v.Set("permission.default_allow", req.DefaultAllow)

	if err := v.WriteConfig(); err != nil {
		response.ServerError(c)
		return
	}

	response.Success(c, map[string]interface{}{
		"message":       "权限系统设置已更新",
		"enabled":       req.Enabled,
		"default_allow": req.DefaultAllow,
	})
}
