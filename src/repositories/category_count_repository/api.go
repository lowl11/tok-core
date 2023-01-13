package category_count_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"tok-core/src/data/entities"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) Get(categoryCode string) (*entities.PostCategoryCountGet, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	ctx, cancel := repo.Ctx()
	defer cancel()

	cursor, err := repo.connection.Find(ctx, mongo_service.Filter().Eq("category_code", categoryCode).Get())
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.PostCategoryCountGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) Create(categoryCode string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.PostCategoryCountCreate{
		CategoryCode: categoryCode,
		Count:        1,
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Increment(categoryCode string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("category_code", categoryCode).Get()
	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"count": 1,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Decrement(categoryCode string) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("category_code", categoryCode).Get()
	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"count": -1,
		},
	}); err != nil {
		return err
	}

	return nil
}
