package mongo

import (
	"context"

	"github.com/sdslabs/katana/types"
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
