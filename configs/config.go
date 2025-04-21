package configs

import (
	"time"

	"github.com/limitcool/starter/internal/pkg/types"
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
	Redis    RedisConfig         // Redis配置
	Log      logconfig.LogConfig // 使用 pkg/logconfig 中的 LogConfig
	Casbin   Casbin              // 权限系统配置
	Storage  Storage             // 文件存储配置
	Admin    Admin               // 管理员配置
	I18n     I18n                // 国际化配置
	GRPC     GRPC                // gRPC服务配置
}

// Config app config
type App struct {
	Port           int
	Name           string
	Mode           string        // 运行模式: debug, release, test
	ReadTimeout    time.Duration // 读取超时
	WriteTimeout   time.Duration // 写入超时
	IdleTimeout    time.Duration // 空闲超时
	MaxHeaderBytes int           // 最大请求头大小
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
	SSLMode         string        // SSL模式，默认disable，可选值：disable, require, verify-ca, verify-full
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

// RedisConfig Redis配置
type RedisConfig struct {
	// 基本配置
	Instances map[string]RedisInstance `yaml:"instances" json:"instances"` // Redis实例配置

	// 缓存配置
	Cache CacheConfig `yaml:"cache" json:"cache"` // 缓存配置
}

// RedisInstance Redis实例配置
type RedisInstance struct {
	Enabled      bool          `yaml:"enabled" json:"enabled"`             // 是否启用
	Addr         string        `yaml:"addr" json:"addr"`                   // 地址
	Password     string        `yaml:"password" json:"password"`           // 密码
	DB           int           `yaml:"db" json:"db"`                       // 数据库索引
	MinIdleConn  int           `yaml:"min_idle_conn" json:"min_idle_conn"` // 最小空闲连接数
	DialTimeout  time.Duration `yaml:"dial_timeout" json:"dial_timeout"`   // 连接超时
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`   // 读取超时
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"` // 写入超时
	PoolSize     int           `yaml:"pool_size" json:"pool_size"`         // 连接池大小
	PoolTimeout  time.Duration `yaml:"pool_timeout" json:"pool_timeout"`   // 连接池超时
	EnableTrace  bool          `yaml:"enable_trace" json:"enable_trace"`   // 是否启用链路追踪
}

// CacheConfig 缓存配置
type CacheConfig struct {
	DefaultTTL        time.Duration `yaml:"default_ttl" json:"default_ttl"`               // 默认过期时间
	KeyPrefix         string        `yaml:"key_prefix" json:"key_prefix"`                 // 键前缀
	EnablePrewarm     bool          `yaml:"enable_prewarm" json:"enable_prewarm"`         // 是否启用预热
	EnableProtection  bool          `yaml:"enable_protection" json:"enable_protection"`   // 是否启用穿透保护
	ProtectionTimeout time.Duration `yaml:"protection_timeout" json:"protection_timeout"` // 穿透保护锁超时
	NilValueTTL       time.Duration `yaml:"nil_value_ttl" json:"nil_value_ttl"`           // 空值缓存过期时间
	LocalCache        bool          `yaml:"local_cache" json:"local_cache"`               // 是否启用本地缓存
	LocalCacheTTL     time.Duration `yaml:"local_cache_ttl" json:"local_cache_ttl"`       // 本地缓存过期时间
	LocalCacheSize    int           `yaml:"local_cache_size" json:"local_cache_size"`     // 本地缓存大小
}

// Casbin 权限系统配置
type Casbin struct {
	Enabled          bool   // 是否启用权限系统
	DefaultAllow     bool   // 默认是否允许访问（权限系统关闭时使用）
	ModelPath        string // Casbin模型文件路径
	PolicyTable      string // 策略表名
	AutoLoad         bool   // 是否自动加载策略
	AutoLoadInterval int    // 自动加载策略间隔（秒）
	LogEnabled       bool   // 是否启用日志
}

// Storage 文件存储配置
type Storage struct {
	Enabled    bool              // 是否启用文件存储
	Type       types.StorageType // 存储类型: local, s3, oss
	Local      LocalStorage      // 本地存储配置
	S3         S3Storage         // S3存储配置
	OSS        OSSStorage        // 阿里云OSS存储配置
	PathConfig PathConfig        // 路径配置
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
	// 用户模式配置
	UserMode string // 用户模式: separate(分离模式), unified(合并模式)
}

// I18n 国际化配置
type I18n struct {
	Enabled          bool     // 是否启用国际化
	DefaultLanguage  string   // 默认语言
	SupportLanguages []string // 支持的语言列表
	ResourcesPath    string   // 语言资源文件路径
}

// GRPC gRPC服务配置
type GRPC struct {
	Enabled      bool          // 是否启用gRPC服务
	Port         int           // gRPC服务端口
	Reflection   bool          // 是否启用反射服务
	HealthCheck  bool          // 是否启用健康检查
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
}
