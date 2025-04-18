package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/limitcool/starter/internal/pkg/storage"
	"github.com/spf13/viper"
)

// LoadConfig 加载配置
func LoadConfig() *Config {
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
		Redis: map[string]Redis{
			"default": {
				Enabled:      false,
				Addr:         "localhost:6379",
				Password:     "",
				DB:           0,
				MinIdleConn:  5,
				DialTimeout:  5,
				ReadTimeout:  3,
				WriteTimeout: 3,
				PoolSize:     10,
				PoolTimeout:  5,
				EnableTrace:  false,
			},
		},
		Log: LogConfig{
			Level:  "info",
			Output: []string{"console"},
			Format: LogFormatText,
			FileConfig: FileLogConfig{
				Path:       "logs",
				MaxSize:    100,
				MaxAge:     30,
				MaxBackups: 10,
				Compress:   false,
			},
			StackTraceEnabled: true,
			StackTraceLevel:   "error",
			MaxStackFrames:    64,
		},
		Casbin: Casbin{
			Enabled:          true,
			DefaultAllow:     false,
			ModelPath:        "configs/rbac_model.conf",
			PolicyTable:      "casbin_rule",
			AutoLoadInterval: 30,
		},
		Storage: Storage{
			Enabled: true,
			Type:    storage.StorageTypeLocal,
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

	// 设置配置文件路径
	configPath := "configs/config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Warn("配置文件不存在，使用默认配置", "path", configPath)
		return config
	}

	// 初始化viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType(filepath.Ext(configPath)[1:])

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		log.Error("读取配置文件失败", "error", err)
		return config
	}

	// 将配置文件内容解析到结构体
	if err := v.Unmarshal(config); err != nil {
		log.Error("解析配置文件失败", "error", err)
		return config
	}

	log.Info("配置加载成功", "path", v.ConfigFileUsed())

	// 打印配置信息
	printConfig(config)

	return config
}

// printConfig 打印配置信息
func printConfig(config *Config) {
	log.Info("应用配置", "name", config.App.Name, "port", config.App.Port)
	log.Info("数据库配置", "enabled", config.Database.Enabled, "driver", config.Driver, "host", config.Database.Host, "port", config.Database.Port)
	log.Info("MongoDB配置", "enabled", config.Mongo.Enabled)
	log.Info("Redis配置", "enabled", config.Redis["default"].Enabled)
	log.Info("Casbin配置", "enabled", config.Casbin.Enabled, "default_allow", config.Casbin.DefaultAllow)
	log.Info("存储配置", "enabled", config.Storage.Enabled, "type", config.Storage.Type)
	log.Info("国际化配置", "enabled", config.I18n.Enabled, "default", config.I18n.DefaultLanguage)
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

	log.Info("配置保存成功", "path", path)
	return nil
}

// structToMap 将结构体转换为map
func structToMap(config *Config) map[string]interface{} {
	v := viper.New()
	v.Set("app", config.App)
	v.Set("driver", config.Driver)
	v.Set("database", config.Database)
	v.Set("jwt_auth", config.JwtAuth)
	v.Set("mongo", config.Mongo)
	v.Set("redis", config.Redis)
	v.Set("log", config.Log)
	v.Set("casbin", config.Casbin)
	v.Set("storage", config.Storage)
	v.Set("admin", config.Admin)
	v.Set("i18n", config.I18n)
	return v.AllSettings()
}
