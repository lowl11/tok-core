package post_like_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"tok-core/src/services/mongo_service"
)

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
