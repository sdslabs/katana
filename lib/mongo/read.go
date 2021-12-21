package mongo

import (
	"context"
	"log"

	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FetchDocs(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []bson.M {
	collection := link.Collection(collectionName)
	var data []bson.M

	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		if err := cur.Decode(&result); err != nil {
			log.Println(err.Error())
			return nil
		}
		data = append(data, result)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return data
}

func FetchSingleTeam(teamName string) (*types.CTFTeam, error) {
	collection := link.Collection(TeamsCollection)
	team := &types.CTFTeam{}
	ctx := context.Background()
	if err := collection.FindOne(ctx, bson.M{UsernameKey: teamName}).Decode(team); err != nil {
		return nil, err
	}
	return team, nil
}

func FetchFlag(flagValue string) (*types.Flag, error) {
	collection := link.Collection(FlagsCollection)
	flag := &types.Flag{}
	ctx := context.Background()
	if err := collection.FindOne(ctx, bson.M{ValueKey: flagValue}).Decode(flag); err != nil {
		return nil, err
	}
	return flag, nil
}

func FetchChallenge(challengeId int) (*types.Challenge, error) {
	collection := link.Collection(ChallengesCollection)
	challenge := &types.Challenge{}
	ctx := context.Background()
	if err := collection.FindOne(ctx, bson.M{IDKey: challengeId}).Decode(challenge); err != nil {
		return nil, err
	}
	return challenge, nil
}

func FetchTeams(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []types.CTFTeam {
	collection := link.Collection(collectionName)
	var data []types.CTFTeam

	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result types.CTFTeam
		if err := cur.Decode(&result); err != nil {
			log.Println(err.Error())
			return nil
		}
		data = append(data, result)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return data
}

func FetchChallenges(ctx context.Context, collectionName string, filter bson.M, opts ...*options.FindOptions) []types.Challenge {
	collection := link.Collection(collectionName)
	var data []types.Challenge

	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result types.Challenge
		if err := cur.Decode(&result); err != nil {
			log.Println(err.Error())
			return nil
		}
		data = append(data, result)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return data
}

func FetchSubmission(teamId int, flag string) (*types.Submission, error) {
	collection := link.Collection(SubmissionsCollection)
	submission := &types.Submission{}
	ctx := context.Background()
	if err := collection.FindOne(ctx, bson.M{SubmittedBy: teamId, Flag: flag}).Decode(submission); err != nil {
		return nil, err
	}
	return submission, nil
}
