package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()

func Connect() (*mongo.Database, error) {
	clientOptions := options.Client()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	mongoUri := fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	clientOptions.ApplyURI(mongoUri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Connect(Ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	return client.Database("kawulo"), nil
}
