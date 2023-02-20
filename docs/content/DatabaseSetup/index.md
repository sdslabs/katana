---
title: "Database"
---

# Introduction

Katana uses a mongoDB to store its data. The database is used to store information like challenge data, user data, and more. This page will walk you through the process of setting up a mongoDB database. The database will run in the master namespace.

# Setup

Simply setup the database by changing the config variables in `config.toml` to the following:

```toml
[mongo]
username = "[YOUR USERNAME HERE]"
password = "[YOUR PASSWORD HERE]"
port = "32000"
mongosh_version = "1.6.1"
```

Default yaml files are written in the `manifests` folder for delpoying mongoDB pods in the master namespace during infraset. To deploy the database, you first need to set up the infrastructure with the help of the `/api/v2/admin/infraSet` endpoint. Then you need to hit the `/api/v2/admin/db` endpoint to setup the database.

# Go Code For Database Setup

The following code is responsible for setting up the database.

In `connection.go`:

```Golang
var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+configs.MongoConfig.Username+":"+configs.MongoConfig.Password+"@"+configs.ServicesConfig.ChallengeDeployer.Host+":"+configs.MongoConfig.Port+"/?directConnection=true&appName=mongosh+"+configs.MongoConfig.Version))
```

In `db.go`:

```Golang
func DB(c *fiber.Ctx) error {
	client, err := utils.GetKubeClient()
	if err != nil {
		log.Println(err)
	}
	service, err := client.CoreV1().Services("default").Get(context.TODO(), "mongo-nodeport-svc", metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	// Print the IP address of the service
	fmt.Println(service.Spec.ClusterIP)
	mongo.Init()

	return c.SendString("Database setup completed")
}
```
