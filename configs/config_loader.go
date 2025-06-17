package configs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/limitcool/starter/internal/pkg/logger"
	"github.com/limitcool/starter/internal/pkg/types"
	"github.com/limitcool/starter/pkg/logconfig"
	"github.com/spf13/viper"
)

// LoadConfig 加载配置
func LoadConfig(configPath string) *Config {
	// 创建默认配置
	config := &Config{
		App: App{
			Port: 8080,
			Name: "Starter",
		},
		Driver: DriverMysql,
		Database: Database{
			Enabled:         true,
			UserName:        "root",
			Password:        "root",
			DBName:          "starter",
			Host:            "localhost",
			Port:            3306,
			TablePrefix:     "",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			ShowLog:         true,
			MaxIdleConn:     10,
			MaxOpenConn:     100,
			ConnMaxLifeTime: 3600,
			SlowThreshold:   500,
		},
		JwtAuth: JwtAuth{
			AccessSecret:  "access_secret",
			AccessExpire:  86400,
			RefreshSecret: "refresh_secret",
			RefreshExpire: 604800,
		},
		Mongo: Mongo{
			Enabled: false,
			URI:     "mongodb://localhost:27017",
			User:    "",
			DB:      "starter",
		},
		Redis: RedisConfig{
			Instances: map[string]RedisInstance{
				"default": {
					Enabled:      true,
					Addr:         "localhost:6379",
					Password:     "",
					DB:           0,
					MinIdleConn:  10,
					DialTimeout:  5 * time.Second,
					ReadTimeout:  3 * time.Second,
					WriteTimeout: 3 * time.Second,
					PoolSize:     100,
					PoolTimeout:  4 * time.Second,
					EnableTrace:  false,
				},
				"session": {
					Enabled:      false,
					Addr:         "localhost:6379",
					Password:     "",
					DB:           1,
					MinIdleConn:  5,
					DialTimeout:  5 * time.Second,
					ReadTimeout:  3 * time.Second,
					WriteTimeout: 3 * time.Second,
					PoolSize:     50,
					PoolTimeout:  4 * time.Second,
					EnableTrace:  false,
				},
			},
			Cache: CacheConfig{
				DefaultTTL:        30 * time.Minute,
				KeyPrefix:         "cache:",
				EnablePrewarm:     true,
				EnableProtection:  true,
				ProtectionTimeout: 5 * time.Second,
				NilValueTTL:       1 * time.Minute,
				LocalCache:        false,
				LocalCacheTTL:     5 * time.Minute,
				LocalCacheSize:    10000,
			},
		},
		Log: logconfig.DefaultLogConfig(),
		Storage: Storage{
			Enabled: true,
			Type:    types.StorageTypeLocal,
			Local: LocalStorage{
				Path: "storage",
				URL:  "/static",
			},
			S3: S3Storage{
				AccessKey: "",
				SecretKey: "",
				Region:    "",
				Bucket:    "",
				Endpoint:  "",
			},
			PathConfig: PathConfig{
				Avatar:    "avatars",
				Document:  "documents",
				Image:     "images",
				Video:     "videos",
				Audio:     "audios",
				Temporary: "temp",
			},
		},
		Admin: Admin{
			Username: "admin",
			Password: "admin123",
			Nickname: "管理员",
		},
		I18n: I18n{
			Enabled:          true,
			DefaultLanguage:  "zh-CN",
			SupportLanguages: []string{"zh-CN", "en-US"},
			ResourcesPath:    "locales",
		},
	}

	// 如果未指定配置文件路径，使用默认路径
	if configPath == "" {
		configPath = "configs/config.yaml"
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			configPath = envPath
		}
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 使用标准库日志，因为此时logger可能还未初始化
		fmt.Printf("配置文件不存在，使用默认配置: %s\n", configPath)
		return config
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType(filepath.Ext(configPath)[1:])

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return config
	}

	// 将配置文件内容解析到结构体
	if err := v.Unmarshal(config); err != nil {
		fmt.Printf("解析配置文件失败: %v\n", err)
		return config
	}

	fmt.Printf("配置加载成功: %s\n", v.ConfigFileUsed())

	return config
}

// PrintConfig 打印配置信息
func PrintConfig(config *Config) {
	// 使用我们的统一logger
	ctx := context.Background()
	logger.InfoContext(ctx, "Application configuration", "app_name", config.App.Name, "app_port", config.App.Port)
	logger.InfoContext(ctx, "Database configuration", "db_enabled", config.Database.Enabled, "db_driver", config.Driver, "db_host", config.Database.Host, "db_port", config.Database.Port)
	logger.InfoContext(ctx, "MongoDB configuration", "mongo_enabled", config.Mongo.Enabled)
	logger.InfoContext(ctx, "Redis configuration", "redis_instances", len(config.Redis.Instances), "redis_default_enabled", config.Redis.Instances["default"].Enabled)
	logger.InfoContext(ctx, "Storage configuration", "storage_enabled", config.Storage.Enabled, "storage_type", config.Storage.Type)
	logger.InfoContext(ctx, "I18n configuration", "i18n_enabled", config.I18n.Enabled, "i18n_default", config.I18n.DefaultLanguage)
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	// 如果路径为空，使用默认路径
	if path == "" {
		path = "configs/config.yaml"
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(filepath.Ext(path)[1:])

	// 将配置结构体转换为map
	if err := v.MergeConfigMap(structToMap(config)); err != nil {
		return fmt.Errorf("合并配置失败: %w", err)
	}

	// 写入文件
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	ctx := context.Background()
	logger.InfoContext(ctx, "Configuration saved successfully", "config_path", path)
	return nil
}

// structToMap 将结构体转换为map
func structToMap(config *Config) map[string]any {
	v := viper.New()
	v.Set("app", config.App)
	v.Set("driver", config.Driver)
	v.Set("database", config.Database)
	v.Set("jwt_auth", config.JwtAuth)
	v.Set("mongo", config.Mongo)
	v.Set("redis", config.Redis)
	v.Set("log", config.Log)
	v.Set("storage", config.Storage)
	v.Set("admin", config.Admin)
	v.Set("i18n", config.I18n)
	return v.AllSettings()
}
