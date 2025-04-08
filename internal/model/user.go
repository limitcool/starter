package model

import (
	"context"

	"github.com/limitcool/starter/internal/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 使用辅助方法获取集合
func getUserCollection() *mongo.Collection {
	return mongodb.Collection("user")
}

type User struct {
}

func (User) Registry() {
	var ctx = context.Background()
	coll := getUserCollection()
	if coll == nil {
		return
	}
	coll.FindOne(ctx, bson.M{"name": "cool"})
	coll.InsertOne(ctx, bson.M{"name": "cool"})
}
