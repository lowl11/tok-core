package category_count_repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"tok-core/src/repositories/repository_mongo"
)

type Repository struct {
	repository_mongo.Base
	connection *mongo.Collection

	mutex sync.Mutex
}

func Create(connection *mongo.Database) *Repository {
	return &Repository{
		Base:       repository_mongo.CreateBase(),
		connection: connection.Collection("category_counts"),

		mutex: sync.Mutex{},
	}
}
