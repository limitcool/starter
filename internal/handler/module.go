package handler

import (
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/filestore"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module handler模块
var Module = fx.Options(
	// 提供所有handler
	fx.Provide(
		func(db *gorm.DB, config *configs.Config, lc fx.Lifecycle) *UserHandler {
			return NewUserHandler(db, config, lc)
		},
		func(db *gorm.DB, config *configs.Config, lc fx.Lifecycle, storage *filestore.Storage) *FileHandler {
			return NewFileHandler(db, config, lc, storage)
		},
		func(db *gorm.DB, config *configs.Config, lc fx.Lifecycle) *AdminHandler {
			return NewAdminHandler(db, config, lc)
		},
	),
)
