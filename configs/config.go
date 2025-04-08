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

// LogFormat 日志格式类型
type LogFormat string

const (
	// LogFormatText 普通文本格式
	LogFormatText LogFormat = "text"
	// LogFormatJSON 结构化JSON格式
	LogFormatJSON LogFormat = "json"
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
	Level      string        // 日志级别: debug, info, warn, error
	Output     []string      // 输出方式: console, file
	Format     LogFormat     // 日志格式: text, json
	FileConfig FileLogConfig // 文件输出配置
}

// FileLogConfig 文件日志配置
type FileLogConfig struct {
	Path       string // 日志文件路径
	MaxSize    int    // 每个日志文件的最大大小（MB）
	MaxAge     int    // 日志文件保留天数
	MaxBackups int    // 保留的旧日志文件最大数量
	Compress   bool   // 是否压缩旧日志文件
}
