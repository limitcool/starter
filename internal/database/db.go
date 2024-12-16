package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/limitcool/starter/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/charmbracelet/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	// 打印所有SQL
	if c.Database.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	// 只打印慢查询
	if c.Database.SlowThreshold > 0 {
		config.Logger = logger.New(
			//将标准输出作为Writer
			log.New(os.Stdout),
			logger.Config{
				//设定慢查询时间阈值
				SlowThreshold: c.Database.SlowThreshold, // nolint: golint
				Colorful:      true,
				//设置日志级别，只有指定级别以上会输出慢查询日志
				LogLevel: logger.Warn,
			},
		)
	}
	config.SkipDefaultTransaction = true
	config.NamingStrategy = schema.NamingStrategy{
		// 使用表前缀
		TablePrefix: c.Database.TablePrefix,
		// 是否使用单数表名
		// SingularTable: true,
	}
	return config
}

func NewDB(c configs.Config) *gorm.DB {
	DB = newDbConn(c)
	return DB
}

// func NewMysql(dsn string, c configs.Config) *gorm.DB {
// 	c.Driver = configs.DriverMysql
// 	return newDbConn(c)
// }

// func NewSqlite(dsn string, c configs.Config) *gorm.DB {
// 	c.Driver = configs.DriverSqlite
// 	return newDbConn(c)
// }

func newDbConn(c configs.Config) *gorm.DB {
	dsn := getDSN(&c)
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
		log.Fatalf("Unsupported database driver: %s", c.Driver)
	}
	if err != nil {
		log.Fatalf("open %s failed. database name: %s, err: %+v", c.Driver, c.Database.DBName, err)
	}
	if c.Driver != configs.DriverSqlite {
		sqlDB.SetMaxOpenConns(c.Database.MaxOpenConn)
		sqlDB.SetMaxIdleConns(c.Database.MaxIdleConn)
		sqlDB.SetConnMaxLifetime(c.Database.ConnMaxLifeTime)
	}

	db, err := gorm.Open(getGormDriver(&c), gormConfig(&c))
	if err != nil {
		log.Fatalf("database connection failed. database name: %s, err: %+v", c.Database.DBName, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.AutoMigrate()
	if err != nil {
		log.Fatal("AutoMigrate err =", err)
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
		log.Fatalf("Unsupported database driver: %s", c.Driver)
		return nil
	}
}
