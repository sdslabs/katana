package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sdslabs/katana/configs"
	"github.com/sdslabs/katana/lib/utils"
	"github.com/sdslabs/katana/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var client *mongo.Client
var link *mongo.Database

func setupAdmin() error {
	adminUser := configs.AdminConfig
	pwd := utils.SHA256(adminUser.Password)
	admin := types.AdminUser{
		Username: adminUser.Username,
		Password: pwd,
	}

	if _, err := AddAdmin(context.Background(), admin); err != nil {
		return fmt.Errorf("cannot add admin: %w", err)
	} else {
		log.Printf("admin privileges have been given to username: %s", admin.Username)
		return nil
	}
}

func setup() error {
	for i := 0; i < 10; i++ {
		log.Printf("Trying to connect to MongoDB, attempt %d", i+1)
		ctx, _ = context.WithTimeout(context.Background(), time.Duration(configs.KatanaConfig.TimeOut)*time.Second)
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+configs.MongoConfig.Username+":"+configs.MongoConfig.Password+"@"+utils.GetKatanaLoadbalancer()+":27017/?directConnection=true&appName=mongosh+"+configs.MongoConfig.Version))
		if err != nil {
			return fmt.Errorf("cannot connect to mongo: %w", err)
		}
		link = client.Database(projectDatabase)
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Println("MongoDB connection was not established")
			log.Println("Error: ", err)
			time.Sleep(time.Duration(configs.KatanaConfig.TimeOut) * time.Second)
		} else {
			log.Println("MongoDB Connection Established")
			if err := setupAdmin(); err != nil {
				return fmt.Errorf("cannot setup admin: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("cannot connect to mongo")
}

func Test() {
	collection := link.Collection(TeamsCollection)
	res, err := collection.InsertOne(context.Background(), bson.M{"team": "ctfteam-0", "password": "passloll"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}

func Init() error {
	if err := setup(); err != nil {
		return err
	}
	return nil
}
