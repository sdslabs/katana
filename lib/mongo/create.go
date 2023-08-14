package mongo

import (
	"context"

	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertOne(ctx context.Context, collectionName string, data interface{}) (interface{}, error) {
	collection := link.Collection(collectionName)
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func InsertMany(ctx context.Context, collectionName string, data []interface{}) ([]interface{}, error) {
	collection := link.Collection(collectionName)
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

func CreateTeams(teams []interface{}) (interface{}, error) {
	return InsertMany(context.Background(), TeamsCollection, teams)
}

func AddAdmin(ctx context.Context, admin types.AdminUser) (interface{}, error) {
	return InsertOne(ctx, AdminCollection, admin)
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
