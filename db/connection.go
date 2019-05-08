package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetDatabase(host, port, user, password, databaseName string, timeout time.Duration) (*mongo.Database, error) {
	connectionString := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		user, password, host, port,
	)

	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(connectionString),
	)
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(
		context.Background(),
		timeout,
	)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(databaseName), nil
}
