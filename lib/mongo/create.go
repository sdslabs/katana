package mongo

import (
	"context"
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
