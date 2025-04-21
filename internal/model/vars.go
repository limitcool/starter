package model

import (
	"time"

	"github.com/limitcool/starter/internal/pkg/idgen"
	"gorm.io/gorm"
)

// BaseModel 基础模型结构 - 普通自增ID
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SnowflakeModel 使用雪花ID的基础模型结构 - 适用于需要分布式ID的模型（如User和SysUser）
type SnowflakeModel struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UUIDModel 使用UUID的基础模型结构 - 适用于需要全局唯一ID且无序的场景
type UUIDModel struct {
	ID        string         `gorm:"primarykey;type:varchar(36)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建记录前自动生成雪花ID
func (m *SnowflakeModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = idgen.GenerateSnowflakeID()
	}
	return nil
}

// BeforeCreate 在创建记录前自动生成UUID
func (m *UUIDModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = idgen.GenerateUUID()
	}
	return nil
}
