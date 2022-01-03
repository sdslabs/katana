package flaghandlerservice

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/mongo"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
)

func Ticker(wg sync.WaitGroup) {
	log.Println("Ticker")
	defer wg.Done()
	handler()
	gocron.Every(uint64(configs.FlagConfig.TickPeriod)).Second().Do(handler)
	<-gocron.Start()
	// gocron.Every(uint64(configs.FlagConfig.TickPeriod)).Second().Do(flaggetter)    todo:trigger setter here
	// <-gocron.Start()
}

func handler() {
	log.Print("Handler triggered -> ", time.Now())
	challenges := mongo.FetchChallenges(context.Background(), mongo.ChallengesCollection, bson.M{})
	teams := mongo.FetchTeams(context.Background(), mongo.TeamsCollection, bson.M{})
	time := time.Now()
	for _, team := range teams {
		for _, challenge := range challenges {
			go flagUpdater(challenge, team, time)
		}
	}
}

func flagUpdater(challenge types.Challenge, team types.CTFTeam, triggerTime time.Time) {
	flagValue := random(configs.FlagConfig.FlagLength)
	var flag = &types.Flag{}
	flag.Value = flagValue
	flag.ChallengeID = challenge.ID
	flag.CreatedAt = triggerTime
	flag.TeamID = team.Index
	if _, err := mongo.InsertOne(context.Background(), mongo.FlagsCollection, flag); err != nil {
		log.Println(err)
	} else {
		log.Println("Trigger setter")
		//todo : trigger flag setter here
	}
}
