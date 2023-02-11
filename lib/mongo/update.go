package mongo

import (
	"context"

	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne(ctx context.Context, collectionName string, filter bson.M, data interface{}, option *options.FindOneAndUpdateOptions) error {
	collection := link.Collection(collectionName)
	return collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": data}, option).Err()
}

func UpsertOne(ctx context.Context, collectionName string, filter bson.M, data interface{}) error {
	return UpdateOne(ctx, collectionName, filter, data, options.FindOneAndUpdate().SetUpsert(true))
}

func AddAdmin(ctx context.Context, admin types.AdminUser) error {
	return UpsertOne(ctx, AdminCollection, bson.M{UsernameKey: admin.Username}, admin)
}

func UpdateFlag(ctx context.Context, flag types.Flag) error {
	return UpsertOne(ctx, FlagsCollection, bson.M{TeamNameKey: flag.Team}, flag)
}
