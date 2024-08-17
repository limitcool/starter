package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/limitcool/starter/configs"
	
	"github.com/limitcool/starter/global"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// RedisClient redis 客户端
var RedisClient *redis.Client

const (
	// ErrRedisNotFound not exist in redis
	ErrRedisNotFound = redis.Nil
	// DefaultRedisName default redis name
	DefaultRedisName = "default"
)

// RedisManager define a redis manager
type RedisManager struct {
	clients map[string]*redis.Client
	config  *configs.Config
	*sync.RWMutex
}

func NewRedisClient(cfg *configs.Config) (*redis.Client, func(), error) {
	clientManager := NewRedisManager(cfg)
	rdb, err := clientManager.GetClient(DefaultRedisName)
	if err != nil {
		return nil, nil, fmt.Errorf("init redis err: %s", err.Error())
	}
	cleanFunc := func() {
		_ = rdb.Close()
	}
	RedisClient = rdb

	return rdb, cleanFunc, nil
}

// NewRedisManager create a redis manager
func NewRedisManager(cfg *configs.Config) *RedisManager {
	return &RedisManager{
		clients: make(map[string]*redis.Client),
		config:  cfg,
		RWMutex: &sync.RWMutex{},
	}
}

// GetClient get a redis instance
func (r *RedisManager) GetClient(name string) (*redis.Client, error) {
	// get client from map
	r.RLock()
	if client, ok := r.clients[name]; ok {
		r.RUnlock()
		return client, nil
	}
	r.RUnlock()

	// create a redis client
	r.Lock()
	defer r.Unlock()

	redisConfig := global.Config.Redis[name]
	fmt.Printf("redisConfig: %v\n", redisConfig)
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		MinIdleConns: redisConfig.MinIdleConn,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
		WriteTimeout: redisConfig.WriteTimeout,
		PoolSize:     redisConfig.PoolSize,
		PoolTimeout:  redisConfig.PoolTimeout,
	})

	// check redis if is ok
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	// hook tracing (using open telemetry)
	if redisConfig.EnableTrace {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			return nil, err
		}
	}
	r.clients[name] = rdb

	return rdb, nil
}
