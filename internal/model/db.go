package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/limitcool/starter/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// NewMySQL 连接数据库，生成数据库实例
func NewMySQL(c *configs.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Mysql.UserName,
		c.Mysql.Password,
		c.Mysql.Host,
		c.Mysql.Port,
		c.Mysql.DBName,
		c.Mysql.Charset,
		c.Mysql.ParseTime,
		//"Asia/Shanghai"),
		c.Mysql.Loc)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("open mysql failed. database name: %s, err: %+v", c.Mysql.DBName, err)
	}
	// set for db connection
	// 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(c.Mysql.MaxOpenConn)
	// 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(c.Mysql.MaxIdleConn)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(c.Mysql.ConnMaxLifeTime)

	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig(c))
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Mysql.DBName, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.AutoMigrate()
	if err != nil {
		log.Fatal("AutoMigrate err =", err)
	}
	DB = db
}

// gormConfig 根据配置决定是否开启日志
func gormConfig(c *configs.Config) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true} // 禁止外键约束, 生产环境不建议使用外键约束
	// 打印所有SQL
	if c.Mysql.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	// 只打印慢查询
	if c.Mysql.SlowThreshold > 0 {
		config.Logger = logger.New(
			//将标准输出作为Writer
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				//设定慢查询时间阈值
				SlowThreshold: c.Mysql.SlowThreshold, // nolint: golint
				Colorful:      true,
				//设置日志级别，只有指定级别以上会输出慢查询日志
				LogLevel: logger.Warn,
			},
		)
	}
	config.SkipDefaultTransaction = true
	config.NamingStrategy = schema.NamingStrategy{
		// 使用表前缀
		TablePrefix: c.Mysql.TablePrefix,
		// 是否使用单数表名
		// SingularTable: true,
	}
	return config
}
