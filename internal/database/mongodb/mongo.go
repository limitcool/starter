package mongodb

import (
	"context"
	"time"

	"github.com/limitcool/starter/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Client

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

// NewMongoDBConn Create new MongoDB client
func NewMongoDBConn(ctx context.Context, cfg *configs.Mongo) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI).
		SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize))
	if err != nil {
		return nil, err
	}
	Mongo = client
	return client, nil
}
