package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetDatabase(host, port, user, password, databaseName string) (*mongo.Database, error) {
	connectionString := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		user, password, host, port,
	)

	ctx, _ := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(connectionString),
	)
	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(
		context.Background(),
		2*time.Second,
	)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(databaseName), nil
}
