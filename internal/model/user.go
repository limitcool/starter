package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var UserCollection = mongo.Collection("user")

type User struct {
}

func (User) Registry() {
	var ctx = context.Background()
	UserCollection.FindOne(ctx, bson.M{"name": "cool"})
	UserCollection.InsertOne(ctx, bson.M{"name": "cool"})
}
