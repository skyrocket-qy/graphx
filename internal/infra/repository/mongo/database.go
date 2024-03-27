package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDb() (*mongo.Client, func(), error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		return nil, nil, err
	}

	var Disconnect = func() {
		client.Disconnect(ctx)
	}
	return client, Disconnect, nil
}
