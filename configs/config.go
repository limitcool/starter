package configs

import (
	"time"
)

type DBDriver string

const (
	DriverMysql    DBDriver = "mysql"
	DriverSqlite   DBDriver = "sqlite3"
	DriverPostgres DBDriver = "postgres"
	DriverMssql    DBDriver = "mssql"
	DriverOracle   DBDriver = "oracle"
	DriverMongo    DBDriver = "mongo"
)

type Config struct {
	App      App
	Driver   DBDriver
	Database Database
	JwtAuth  JwtAuth
	Mongo    Mongo
}

// Config app config
type App struct {
	Port int
	Name string
}

// Config mysql config
type Database struct {
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

// Config MongoDB config
type Mongo struct {
	URI      string
	User     string
	Password string
	DB       string
}
