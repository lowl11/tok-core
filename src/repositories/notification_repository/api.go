package notification_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) GetInfo(username string) ([]entities.NotificationGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Get()

	cursor, err := repo.connection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	list := make([]entities.NotificationGet, 0)
	for cursor.Next(ctx) {
		item := entities.NotificationGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) GetCount(username string) (int, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("status", "new").Get()

	count, err := repo.connection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *Repository) Create(username, actionAuthor, actionKey, actionCode string, body *entities.NotificationBody) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.NotificationCreate{
		Username:     username,
		Status:       "new",
		ActionAuthor: actionAuthor,
		ActionKey:    actionKey,
		ActionCode:   actionCode,
		ActionBody:   body,
		CreatedAt:    time.Now(),
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ReadItem(username, actionKey string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("action_key", actionKey).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"status": "read",
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ReadItemList(username string, actionKeys []string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("username", username).Eq("action_key", bson.M{
		"$in": actionKeys,
	}).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$set": bson.M{
			"status": "read",
		},
	}); err != nil {
		return err
	}

	return nil
}
