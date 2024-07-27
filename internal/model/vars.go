package model

import (
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/database"
	"github.com/limitcool/starter/internal/database/mongodb"
	"gorm.io/gorm"
)

var db = database.DB

var mongo = mongodb.Mongo.Database(global.Config.Mongo.DB)

var ErrNotFound = gorm.ErrRecordNotFound
