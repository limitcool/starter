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
	"github.com/limitcool/starter/common/redis"
	"github.com/limitcool/starter/configs"
	"github.com/limitcool/starter/global"
	"github.com/limitcool/starter/internal/database"
	"github.com/limitcool/starter/internal/database/mongodb"
	"github.com/limitcool/starter/routers"

	"github.com/spf13/viper"
)

func main() {
	lib.SetDebugMode(func() {
		log.Info("Debug Mode")
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	})

	log.SetPrefix("🌏 starter ")
	viper.SetConfigFile("./config.yaml")
	viper.ReadInConfig()
	err := viper.Unmarshal(&global.Config)
	if err != nil {
		log.Fatal("viper unmarshal err = ", err)
	}
	switch global.Config.Driver {
	case "":
		// log.Fatal("driver is empty")
	case configs.DriverMongo:
		log.Info("driver is mongo")
		_, err := mongodb.NewMongoDBConn(context.Background(), &global.Config.Mongo)
		if err != nil {
			log.Fatal("mongo connect err = ", err)
		}
	default:
		log.Info("driver is ", global.Config.Driver)
		db := database.NewDB(*global.Config)
		db.AutoMigrate()
	}
	_, _, err = redis.NewRedisClient(global.Config)
	if err != nil {
		log.Fatal("redis connect err = ", err)
	}
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
