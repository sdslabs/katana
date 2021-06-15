package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sdslabs/katana/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.MongoConfig.URL))
var link = client.Database(projectDatabase)

func Test() {
	collection := link.Collection(teamsCollection)
	res, err := collection.InsertOne(context.Background(), bson.M{"team": "ctfteam-0", "password": "passloll"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
