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

func AddChallenge(challenge types.Challenge, teamName string) error {
	teamFilter := bson.M{"username": teamName}
	update := bson.M{"$push": bson.M{"challenges": challenge}}
	_, err := link.Collection(TeamsCollection).UpdateOne(context.TODO(), teamFilter, update)
	if err != nil {
		return err
	}
	return nil
}
