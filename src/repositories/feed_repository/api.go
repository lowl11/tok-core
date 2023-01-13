package feed_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"tok-core/src/data/entities"
	"tok-core/src/services/mongo_service"
)

const (
	Explore      = "explore"
	Unauthorized = "unauthorized"
)

func (repo *Repository) Get(feedName string) (*entities.FeedGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("name", feedName).Get()

	cursor, err := repo.connection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.FeedGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) Create(feedName string, post *entities.FeedPost) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.FeedCreate{
		Name:  feedName,
		Count: 1,
		Posts: []entities.FeedPost{*post},
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) CreateWithList(feedName string, posts []entities.FeedPost) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := &entities.FeedCreate{
		Name:  feedName,
		Count: len(posts),
		Posts: posts,
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) AddPostExist(feedName string, entity *entities.FeedPost) error {
	feed, err := repo.Get(feedName)
	if err != nil {
		return err
	}

	if feed != nil {
		if err = repo.AddPost(feedName, entity); err != nil {
			return err
		}
	} else {
		if err = repo.Create(feedName, entity); err != nil {
			return err
		}
	}

	return nil
}

func (repo *Repository) AddPost(feedName string, entity *entities.FeedPost) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("name", feedName).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{
			"posts": entity,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) AddPostListExist(feedName string, list []entities.FeedPost) error {
	feed, err := repo.Get(feedName)
	if err != nil {
		return err
	}

	if feed != nil {
		if err = repo.AddPostList(feedName, list); err != nil {
			return err
		}
	} else {
		if err = repo.CreateWithList(feedName, list); err != nil {
			return err
		}
	}

	return nil
}

func (repo *Repository) AddPostList(feedName string, list []entities.FeedPost) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("name", feedName).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{
			"posts": list,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) RemovePost(feedName, postCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("name", feedName).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$pull": bson.M{
			"posts.post_code": postCode,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) RemovePostList(feedName, postCodes []string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("name", feedName).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$pull": bson.M{
			"posts.post_code": postCodes,
		},
	}); err != nil {
		return err
	}

	return nil
}
