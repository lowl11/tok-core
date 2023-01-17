package post_like_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) Get(postCode string) (*entities.PostLikeGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	result := repo.connection.FindOne(ctx, mongo_service.Filter().Eq("post_code", postCode).Get())
	if result.Err() != nil {
		if strings.Contains(result.Err().Error(), "no documents") {
			return nil, nil
		}

		return nil, result.Err()
	}

	item := entities.PostLikeGet{}
	if err := result.Decode(&item); err != nil {
		return nil, err
	}

	if item.PostCode == "" {
		return nil, nil
	}

	return &item, nil
}

func (repo *Repository) GetByList(postCodes []string) ([]entities.PostLikeGetList, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := bson.M{
		"post_code": bson.M{
			"$in": postCodes,
		},
	}

	cursor, err := repo.connection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	list := make([]entities.PostLikeGetList, 0)
	for cursor.Next(ctx) {
		item := entities.PostLikeGetList{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) Create(model *models.PostLike, likeAuthor string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	if _, err := repo.connection.InsertOne(ctx, &entities.PostLikeCreate{
		PostCode:    model.PostCode,
		LikesCount:  1,
		LikeAuthors: []string{likeAuthor},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Like(postCode, likeAuthor string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("post_code", postCode).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{
			"like_authors": likeAuthor,
		},
		"$inc": bson.M{
			"likes_count": 1,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Unlike(postCode, likeAuthor string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("post_code", postCode).Get()

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$pull": bson.M{
			"like_authors": likeAuthor,
		},
		"$inc": bson.M{
			"likes_count": -1,
		},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DeleteByPost(postCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	filter := mongo_service.Filter().Eq("post_code", postCode).Get()
	if _, err := repo.connection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
