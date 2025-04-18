package sqldb

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/limitcool/starter/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/charmbracelet/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量
var (
	dbOnce sync.Once
)

func getDSN(c *configs.Config) string {
	switch c.Driver {
	case configs.DriverSqlite:
		return fmt.Sprintf("%s.db", c.Database.DBName)
	case configs.DriverPostgres:
		return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			c.Database.UserName,
			c.Database.Password,
			c.Database.Host,
			c.Database.DBName,
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
	gormLogger := logger.New(
		// 使用我们的结构化日志
		&gormLogWriter{
			logger: log.Default(),
		},
		logger.Config{
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
type gormLogWriter struct {
	logger *log.Logger
}

// Printf 实现Print接口供GORM日志使用
func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	// 将GORM日志输出到结构化日志
	msg := fmt.Sprintf(format, args...)
	w.logger.Info("GORM", "message", msg)
}

// getGormLogLevel 根据配置获取GORM日志级别
func getGormLogLevel(c *configs.Config) logger.LogLevel {
	if c.Database.ShowLog {
		return logger.Info
	}

	if c.Database.SlowThreshold > 0 {
		return logger.Warn
	}

	return logger.Silent
}

// NewDB 创建数据库连接
func NewDB(c configs.Config) *gorm.DB {
	var db *gorm.DB
	dbOnce.Do(func() {
		db = newDbConn(&c)
		// 设置全局实例
		setupInstance(&Component{db: db, Config: &c, enabled: c.Database.Enabled})
	})
	return Instance().DB()
}

// func NewMysql(dsn string, c configs.Config) *gorm.DB {
// 	c.Driver = configs.DriverMysql
// 	return newDbConn(c)
// }

// func NewSqlite(dsn string, c configs.Config) *gorm.DB {
// 	c.Driver = configs.DriverSqlite
// 	return newDbConn(c)
// }

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
		log.Fatal("Unsupported database driver", "driver", c.Driver)
	}
	if err != nil {
		log.Fatal("Failed to open database connection",
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
			log.Warn("Database name is empty, using default", "driver", c.Driver)
			c.Database.DBName = "default"
		}
	}

	db, err := gorm.Open(getGormDriver(c), gormConfig(c))
	if err != nil {
		log.Fatal("Database connection failed",
			"database", c.Database.DBName,
			"error", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.AutoMigrate()
	if err != nil {
		log.Fatal("AutoMigrate failed", "error", err)
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
		log.Fatal("Unsupported database driver", "driver", c.Driver)
		return nil
	}
}
