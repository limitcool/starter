package configs

import (
	"time"
)

type Config struct {
	App     App
	Mysql   Mysql
	JwtAuth JwtAuth
}

// Config app config
type App struct {
	Port int
}

// Config mysql config
type Mysql struct {
	UserName        string
	Password        string
	DBName          string
	Host            string
	Port            int
	TablePrefix     string
	Charset         string
	ParseTime       bool
	Loc             string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime time.Duration
	SlowThreshold   time.Duration // 慢查询时长，默认500ms
}

// Config jwt config
type JwtAuth struct {
	AccessSecret string
	AccessExpire int64
}
