package examples

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongoDB(ctx context.Context) *mongo.Database {
	client, err := mongo.NewClient(
		options.Client().
			ApplyURI("mongodb://localhost:27017").
			SetAuth(options.Credential{
				Username: "root",
				Password: "abcd1234",
			}))
	if err != nil {
		panic(err)
	}
	if err := client.Connect(ctx); err != nil {
		panic(err)
	}
	return client.Database("sqlike")
}
