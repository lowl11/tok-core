package user_interest_repository

import (
	"github.com/lowl11/lazy-collection/array"
	"go.mongodb.org/mongo-driver/bson"
	"tok-core/src/data/entities"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) Get(username string) (*entities.UserInterestGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()
	cursor, err := repo.interest.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.UserInterestGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) Create(username, categoryCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.UserInterestCreate{
		Username: username,
		Categories: []entities.UserInterestCategory{
			{
				CategoryCode: categoryCode,
				Interest:     increaseInterest,
			},
		},
	}

	if _, err := repo.interest.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) IncreaseCategoryExist(username, categoryCode string) error {
	interest, err := repo.Get(username)
	if err != nil {
		return err
	}

	if interest != nil {
		if err = repo.IncreaseCategory(username, categoryCode); err != nil {
			return err
		}
	} else {
		if err = repo.Create(username, categoryCode); err != nil {
			return err
		}
	}

	return nil
}

func (repo *Repository) IncreaseCategory(username, categoryCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("categories.category_code", categoryCode).Get()

	if _, err := repo.interest.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"categories.$.interest": increaseInterest,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DecreaseCategoryExist(username, categoryCode string) error {
	interest, err := repo.Get(username)
	if err != nil {
		return err
	}

	// если записи нет, то и уменьшать не нужно
	if interest == nil {
		return nil
	}

	interestCategory := array.NewWithList[entities.UserInterestCategory](interest.Categories...).Single(func(item entities.UserInterestCategory) bool {
		return item.CategoryCode == categoryCode
	})

	// если категории нет, то и уменьшать не нужно
	if interestCategory == nil || interestCategory.Interest <= 0 {
		return nil
	}

	if err = repo.DecreaseCategory(username, categoryCode); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DecreaseCategory(username, categoryCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("categories.category_code", categoryCode).Get()

	if _, err := repo.interest.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"categories.$.interest": -1 * increaseInterest,
		},
	}); err != nil {
		return err
	}

	return nil
}
