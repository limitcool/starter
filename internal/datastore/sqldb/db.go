package sqldb

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/glebarez/sqlite"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/internal/pkg/logger"
	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getDSN(c *configs.Config) string {
	switch c.Driver {
	case configs.DriverSqlite:
		return fmt.Sprintf("%s.db", c.Database.DBName)
	case configs.DriverPostgres:
		// 使用url.QueryEscape处理密码中的特殊字符
		password := url.QueryEscape(c.Database.Password)

		// 构建连接字符串
		sslMode := "disable"
		if c.Database.SSLMode != "" {
			sslMode = c.Database.SSLMode
		}

		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			c.Database.UserName,
			password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
			sslMode,
		)
	case configs.DriverMysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			c.Database.UserName,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
			c.Database.Charset,
			c.Database.ParseTime,
			c.Database.Loc,
		)
	default:
		c.Driver = configs.DriverSqlite
		return fmt.Sprintf("%s.db", c.Database.DBName)
	}

}

// gormConfig 根据配置决定是否开启日志
func gormConfig(c *configs.Config) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true} // 禁止外键约束, 生产环境不建议使用外键约束

	// 创建一个结构化日志适配器
	gormLogger := gormlogger.New(
		// 使用我们的结构化日志
		&gormLogWriter{},
		gormlogger.Config{
			SlowThreshold: c.Database.SlowThreshold,
			Colorful:      false, // 禁用颜色以避免与结构化日志冲突
			LogLevel:      getGormLogLevel(c),
		},
	)

	// 设置GORM日志
	config.Logger = gormLogger

	config.SkipDefaultTransaction = true
	config.NamingStrategy = schema.NamingStrategy{
		// 使用表前缀
		TablePrefix: c.Database.TablePrefix,
		// 是否使用单数表名
		// SingularTable: true,
	}
	return config
}

// gormLogWriter 适配GORM日志到结构化日志
type gormLogWriter struct{}

// Printf 实现Print接口供GORM日志使用
func (w *gormLogWriter) Printf(format string, args ...any) {
	// 将GORM日志输出到结构化日志
	msg := fmt.Sprintf(format, args...)
	logger.Info("GORM", "message", msg)
}

// getGormLogLevel 根据配置获取GORM日志级别
func getGormLogLevel(c *configs.Config) gormlogger.LogLevel {
	if c.Database.ShowLog {
		return gormlogger.Info
	}

	if c.Database.SlowThreshold > 0 {
		return gormlogger.Warn
	}

	return gormlogger.Silent
}

// NewDBWithConfig 创建数据库连接
func NewDBWithConfig(c configs.Config) *gorm.DB {
	return newDbConn(&c)
}

func newDbConn(c *configs.Config) *gorm.DB {
	dsn := getDSN(c)
	var (
		err   error
		sqlDB *sql.DB
	)
	switch c.Driver {
	case configs.DriverMysql:
		sqlDB, err = sql.Open("mysql", dsn)
	case configs.DriverPostgres:
		sqlDB, err = sql.Open("pgx", dsn)
	case configs.DriverSqlite:
		// sqlDB, err = sql.Open("sqlite3", dsn) // 注意：SQLite 的驱动名称是 "sqlite3"

	default:
		logger.Fatal("Unsupported database driver", "driver", c.Driver)
	}
	if err != nil {
		logger.Fatal("Failed to open database connection",
			"driver", c.Driver,
			"database", c.Database.DBName,
			"error", err)
	}
	if c.Driver != configs.DriverSqlite {

		sqlDB.SetMaxOpenConns(c.Database.MaxOpenConn)
		sqlDB.SetMaxIdleConns(c.Database.MaxIdleConn)
		sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifeTime)
	} else {
		if c.Database.DBName == "" {
			logger.Warn("Database name is empty, using default", "driver", c.Driver)
			c.Database.DBName = "default"
		}
	}

	db, err := gorm.Open(getGormDriver(c), gormConfig(c))
	if err != nil {
		logger.Fatal("Database connection failed",
			"database", c.Database.DBName,
			"error", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.AutoMigrate()
	if err != nil {
		logger.Fatal("AutoMigrate failed", "error", err)
	}
	return db
}

func getGormDriver(c *configs.Config) gorm.Dialector {
	switch c.Driver {
	case configs.DriverMysql:
		return mysql.Open(getDSN(c))
	case configs.DriverPostgres:
		return postgres.Open(getDSN(c))
	case configs.DriverSqlite:
		return sqlite.Open(getDSN(c))
	default:
		logger.Fatal("Unsupported database driver", "driver", c.Driver)
		return nil
	}
}
