package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateOne(ctx context.Context, collectionName string, filter bson.M, data interface{}, option *options.FindOneAndUpdateOptions) error {
	collection := link.Collection(collectionName)
	return collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": data}, option).Err()
}
