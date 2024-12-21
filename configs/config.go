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
	Redis    map[string]Redis
	Log      LogConfig
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

// Redis配置结构
type Redis struct {
	Addr         string
	Password     string
	DB           int
	MinIdleConn  int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
	EnableTrace  bool
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string   `yaml:"level"`       // 日志级别: debug, info, warn, error
	Output     []string `yaml:"output"`      // 输出方式: console, file
	FileConfig FileLogConfig `yaml:"file"`   // 文件输出配置
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path        string `yaml:"path"`         // 日志文件路径
	MaxSize     int    `yaml:"maxSize"`      // 每个日志文件的最大大小（MB）
	MaxAge      int    `yaml:"maxAge"`       // 日志文件保留天数
	MaxBackups  int    `yaml:"maxBackups"`   // 保留的旧日志文件最大数量
	Compress    bool   `yaml:"compress"`     // 是否压缩旧日志文件
}
