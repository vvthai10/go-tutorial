package bootstrap

import (
	"context"
	"log"
	"time"

	"github.com/vvthai10/transaction-mongodb/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoDB(env *config.Env) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(env.DBUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary()) 
	if err != nil {
		log.Fatal(err)
	}
	
	db := client.Database("transaction-db")

	return db
}