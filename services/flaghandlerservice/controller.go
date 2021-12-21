package flaghandlerservice

import (
	"context"
	"time"

	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkSubmission(submission *types.Submission) bool {
	_, err := mongo.FetchSubmission(submission.SubmittedBy, submission.Flag)
	return err != nil
}

func submitFlag(value string, team types.CTFTeam) (bool, int) {
	flag, err := mongo.FetchFlag(value)
	var points int
	if err != nil {
		return false, 0
	}
	if flag.TeamID == team.Index {
		return false, 0
	}
	submission := &types.Submission{}
	submission.ChallengeID = flag.ChallengeID
	submission.Flag = flag.Value
	submission.SubmittedBy = team.Index
	submission.Time = time.Now()
	if !checkSubmission(submission) {
		return false, 0
	}
	if res, err := mongo.FetchChallenge(flag.ChallengeID); err != nil {
		return false, 0
	} else {
		team.Score = team.Score + res.Points
		points = res.Points
	}
	if err := mongo.UpdateOne(context.Background(), mongo.TeamsCollection, bson.M{"id": team.Index}, team, options.FindOneAndUpdate().SetUpsert(false)); err != nil {
		return false, 0
	}
	if _, err := mongo.InsertOne(context.Background(), mongo.SubmissionsCollection, submission); err != nil {
		return false, 0
	}
	return true, points
}
