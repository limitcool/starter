package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/charmbracelet/log"
	"github.com/limitcool/lib"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/database"
	"github.com/limitcool/starter/internal/database/mongodb"
	"github.com/limitcool/starter/routers"

	"github.com/limitcool/starter/pkg/env"
	"github.com/spf13/viper"
	"github.com/limitcool/starter/pkg/logger"
)

func loadConfig() {
	env := env.Get()
	log.Info("current env:", env)

	// 设置默认配置文件
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")

	// 读取默认配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("read default config err = ", err)
	}

	// 读取环境配置
	viper.SetConfigName(fmt.Sprintf("config-%s", env))
	if err := viper.MergeInConfig(); err != nil {
		log.Warn("config not found, use default config")
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&global.Config); err != nil {
		log.Fatal("config unmarshal err = ", err)
	}
}

func main() {
	lib.SetDebugMode(func() {
		log.Info("Debug Mode")
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	})

	log.SetPrefix("🌏 starter ")

	// 加载配置
	loadConfig()

	// 初始化日志
	logger.Setup(global.Config.Log)

	switch global.Config.Driver {
	case configs.DriverMongo:
		log.Info("driver is mongo")
		_, err := mongodb.NewMongoDBConn(context.Background(), &global.Config.Mongo)
		if err != nil {
			log.Fatal("mongo connect err = ", err)
		}
	case configs.DriverMysql, configs.DriverPostgres, configs.DriverSqlite, configs.DriverMssql, configs.DriverOracle:
		log.Info("driver is ", global.Config.Driver)
		db := database.NewDB(*global.Config)
		db.AutoMigrate()
	default:
		log.Info("driver is none")
	}
	// _, _, err = redis.NewRedisClient(global.Config)
	// if err != nil {
	// 	log.Fatal("redis connect err = ", err)
	// }
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           fmt.Sprint("0.0.0.0:", global.Config.App.Port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}
	log.Infof("Listen: %s:%d\n", "http://127.0.0.1", global.Config.App.Port)
	go func() {
		// 服务连接 监听
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("Listen:%s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器,这里需要缓冲
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	//(设置5秒超时时间)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		// 处理错误，例如记录日志、返回错误等
		log.Infof("Error during shutdown: %v", err)
	}
}
