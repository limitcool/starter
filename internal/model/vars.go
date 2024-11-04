package model

import (
	"time"

	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/database"
	"github.com/limitcool/starter/internal/database/mongodb"
	"gorm.io/gorm"
)

var db = database.DB

var mongo = mongodb.Mongo.Database(global.Config.Mongo.DB)

var ErrNotFound = gorm.ErrRecordNotFound

type BaseModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
