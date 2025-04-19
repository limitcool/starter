package configs

import (
	"time"

	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/limitcool/starter/pkg/logconfig"
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
	Log      logconfig.LogConfig // 使用 pkg/logconfig 中的 LogConfig
	Casbin   Casbin             // 权限系统配置
	Storage  Storage            // 文件存储配置
	Admin    Admin              // 管理员配置
	I18n     I18n               // 国际化配置
}

// Config app config
type App struct {
	Port int
	Name string
}

// Config MySQL等数据库配置
type Database struct {
	Enabled         bool // 是否启用SQL数据库
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
	AccessSecret  string
	AccessExpire  int64
	RefreshSecret string
	RefreshExpire int64
}

// Config MongoDB config
type Mongo struct {
	Enabled  bool // 是否启用MongoDB
	URI      string
	User     string
	Password string
	DB       string
}

// Redis配置结构
type Redis struct {
	Enabled      bool // 是否启用Redis
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

// Casbin 权限系统配置
type Casbin struct {
	Enabled          bool   // 是否启用权限系统
	DefaultAllow     bool   // 默认是否允许访问（权限系统关闭时使用）
	ModelPath        string // Casbin模型文件路径
	PolicyTable      string // 策略表名
	AutoLoadInterval int    // 自动加载策略间隔（秒）
}

// Storage 文件存储配置
type Storage struct {
	Enabled    bool                // 是否启用文件存储
	Type       storage.StorageType // 存储类型: local, s3, oss
	Local      LocalStorage        // 本地存储配置
	S3         S3Storage           // S3存储配置
	OSS        OSSStorage          // 阿里云OSS存储配置
	PathConfig PathConfig          // 路径配置
}

// LocalStorage 本地存储配置
type LocalStorage struct {
	Path string // 本地存储路径
	URL  string // 访问URL前缀
}

// S3Storage AWS S3存储配置
type S3Storage struct {
	AccessKey string // 访问密钥ID
	SecretKey string // 访问密钥Secret
	Region    string // 区域
	Bucket    string // 桶名称
	Endpoint  string // 端点URL
}

// OSSStorage 阿里云OSS存储配置
type OSSStorage struct {
	AccessKey string // 访问密钥ID
	SecretKey string // 访问密钥Secret
	Region    string // 区域
	Bucket    string // 桶名称
	Endpoint  string // 端点URL
}

// PathConfig 存储路径配置
type PathConfig struct {
	Avatar    string // 头像存储路径
	Document  string // 文档存储路径
	Image     string // 图片存储路径
	Video     string // 视频存储路径
	Audio     string // 音频存储路径
	Temporary string // 临时文件存储路径
}

// Admin 管理员配置
type Admin struct {
	Username string // 管理员用户名
	Password string // 管理员密码
	Nickname string // 管理员昵称
}

// I18n 国际化配置
type I18n struct {
	Enabled          bool     // 是否启用国际化
	DefaultLanguage  string   // 默认语言
	SupportLanguages []string // 支持的语言列表
	ResourcesPath    string   // 语言资源文件路径
}
