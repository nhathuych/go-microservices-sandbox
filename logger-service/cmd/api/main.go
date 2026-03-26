package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	serverPort = "80"
	rpcPort    = "5001"
	gRpcPort   = "50051"
	mongoURL   = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct{}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}

	if err := c.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return c, nil
}
