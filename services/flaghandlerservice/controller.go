package flaghandlerservice

import (
	"context"
	"log"

	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func submitFlag(value string, team types.CTFTeam) (bool, int) {
	flag, err := mongo.FetchFlag(value)
	var points int
	if err != nil {
		if flag.TeamID == team.Index {
			return false, 0
		}
		submission := &types.Submission{}
		submission.ChallengeID = flag.ChallengeID
		submission.Flag = flag.Value
		submission.Submitter = team.Index //time of submission could also be stored
		if res, err := mongo.InsertOne(context.Background(), mongo.SubmissionsCollection, submission); err != nil {
			log.Println(err)
			return false, 0
		} else {
			log.Println(res)
		}
		if res, err := mongo.FetchChallenge(flag.ChallengeID); err != nil {
			log.Println(err)
			return false, 0
		} else {
			team.Score += res.Points
			points = res.Points
		}
		if err := mongo.UpdateOne(context.Background(), mongo.TeamsCollection, bson.M{"id": team.Index}, team, options.FindOneAndUpdate().SetUpsert(false)); err != nil {
			log.Println(err)
			return false, 0
		}
		return true, points
	}
	return false, 0
}
