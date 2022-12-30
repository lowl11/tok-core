package mongo_service

import (
	"context"
	"github.com/lowl11/lazylog/layers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"tok-core/src/definition"
)

func NewConnection(databaseName string) (*mongo.Database, error) {
	config := definition.Config.Mongo
	logger := definition.Logger

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Connection))
	if err != nil {
		return nil, err
	}

	logger.Info("Ping Mongo database...", layers.Mongo)
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	logger.Info("Ping Mongo database is done!", layers.Mongo)

	return client.Database(databaseName), nil
}
