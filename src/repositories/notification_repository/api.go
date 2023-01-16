package notification_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"tok-core/src/data/entities"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) GetInfo(username string) (*entities.NotificationGetInfo, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()

	cursor, err := repo.connection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.NotificationGetInfo{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) GetCount(username string) (*entities.NotificationGetCount, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()

	cursor, err := repo.connection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.NotificationGetCount{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) Create(username string, item *entities.NotificationAction) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.NotificationCreate{
		Username:        username,
		NewActionsCount: 1,
		Actions:         []entities.NotificationAction{*item},
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) AddItem(username string, item *entities.NotificationAction) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{
			"actions": item,
		},
		"$inc": bson.M{
			"new_actions_count": 1,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ReadItem(username, actionKey string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("actions.action_key", actionKey).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"new_actions_count": -1,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ReadItemList(username string, actionKeys []string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{
			"new_actions_count": -1 * len(actionKeys),
		},
	}); err != nil {
		return err
	}

	return nil
}
