package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/limitcool/starter/internal/api/response"
	"github.com/limitcool/starter/internal/core"
)

// 获取系统设置

func NewSystemController() *SystemController {
	return &SystemController{}
}

type SystemController struct {
}

func (sc *SystemController) GetSystemSettings(c *gin.Context) {
	// 获取配置
	config := core.Instance().Config()

	// 返回当前权限系统设置
	settings := map[string]interface{}{
		"permission": map[string]interface{}{
			"enabled":       config.Permission.Enabled,
			"default_allow": config.Permission.DefaultAllow,
		},
	}

	response.Success(c, settings)
}
