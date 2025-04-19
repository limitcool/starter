package migration

import (
	"github.com/limitcool/starter/internal/model"
	"gorm.io/gorm"
)

// 这个文件用于初始化 Casbin 权限规则

func init() {
	// 注册 Casbin 权限规则迁移
	RegisterMigration("017_init_casbin_rules", initCasbinRules, dropCasbinRules)
}

// CasbinRule Casbin规则模型
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:100"`
	V0    string `gorm:"size:100"`
	V1    string `gorm:"size:100"`
	V2    string `gorm:"size:100"`
	V3    string `gorm:"size:100"`
	V4    string `gorm:"size:100"`
	V5    string `gorm:"size:100"`
}

// TableName 表名
func (CasbinRule) TableName() string {
	return "casbin_rule"
}

// initCasbinRules 初始化 Casbin 权限规则
func initCasbinRules(tx *gorm.DB) error {
	// 检查是否已有规则
	var count int64
	if err := tx.Model(&CasbinRule{}).Count(&count).Error; err != nil {
		return err
	}

	// 已存在则不重复创建
	if count > 0 {
		return nil
	}

	// 获取管理员角色
	var adminRole model.Role
	if err := tx.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// 为管理员角色添加所有权限
	rules := []CasbinRule{
		// 角色继承关系
		{Ptype: "g", V0: "1", V1: "admin"}, // 用户ID为1的用户继承admin角色的所有权限

		// 管理员角色拥有所有权限
		{Ptype: "p", V0: "admin", V1: "*", V2: "*"}, // admin角色可以访问任何路径的任何方法
	}

	// 创建规则
	for _, rule := range rules {
		if err := tx.Create(&rule).Error; err != nil {
			return err
		}
	}

	return nil
}

// dropCasbinRules 删除 Casbin 权限规则
func dropCasbinRules(tx *gorm.DB) error {
	return tx.Delete(&CasbinRule{}).Error
}
