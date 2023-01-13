package feed_repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"tok-core/src/repositories/repository_mongo"
)

type Repository struct {
	repository_mongo.Base
	connection *mongo.Collection
}

func Create(connection *mongo.Database) *Repository {
	return &Repository{
		Base:       repository_mongo.CreateBase(),
		connection: connection.Collection("feeds"),
	}
}
