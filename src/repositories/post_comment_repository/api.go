package post_comment_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) GetByCode(commentCode string, subComment bool) (*entities.PostCommentGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	var filter bson.M
	if subComment {
		filter = mongo_service.Filter().Eq("comments.subcomments.comment_code", commentCode).Get()
	} else {
		filter = mongo_service.Filter().Eq("comments.comment_code", commentCode).Get()
	}

	result := repo.connection.FindOne(ctx, filter)
	if result.Err() != nil {
		if strings.Contains(result.Err().Error(), "no documents") {
			return nil, nil
		}
		return nil, result.Err()
	}

	item := entities.PostCommentGet{}
	if err := result.Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *Repository) GetByPost(postCode string) (*entities.PostCommentGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	cursor, err := repo.connection.Find(ctx, mongo_service.Filter().Eq("post_code", postCode).Get())
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)
	defer repo.LogError(cursor.Err())

	if cursor.Next(ctx) {
		item := entities.PostCommentGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}

func (repo *Repository) GetByList(postCodes []string) ([]entities.PostCommentGetList, error) {
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

	list := make([]entities.PostCommentGetList, 0)
	for cursor.Next(ctx) {
		item := entities.PostCommentGetList{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) Create(model *models.PostCommentAdd, commentAuthor, commentCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := entities.PostCommentCreate{
		PostCode:   model.PostCode,
		PostAuthor: model.PostAuthor,

		CommentsCount: 1,
		Comments: []entities.PostCommentItem{
			{
				CommentCode:   commentCode,
				CommentAuthor: commentAuthor,
				CommentText:   model.CommentText,
				CreatedAt:     time.Now(),
				SubComments:   []entities.PostSubCommentItem{},
			},
		},
	}

	if _, err := repo.connection.InsertOne(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Append(model *models.PostCommentAdd, commentAuthor, commentCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// определяем "подкоммент" ли это
	// определяется это тем что фронт отправляет родительский код коммента
	// если он есть, значит это "подкоммент"
	isSubComment := model.ParentCommentCode != ""

	var entity any
	var filter bson.M
	var pushName string

	createdAt := time.Now()
	likeAuthors := make([]string, 0)

	if isSubComment {
		entity = entities.PostCommentAppendSubComment{
			CommentCode:   commentCode,
			CommentAuthor: commentAuthor,
			CommentText:   model.CommentText,
			LikeAuthors:   likeAuthors,
			CreatedAt:     createdAt,
		}

		filter = mongo_service.Filter().Eq("post_code", model.PostCode).Eq("comments.comment_code", model.ParentCommentCode).Get()
		pushName = "comments.$.subcomments"
	} else {
		entity = entities.PostCommentAppendComment{
			CommentCode:   commentCode,
			CommentAuthor: commentAuthor,
			CommentText:   model.CommentText,
			LikeAuthors:   likeAuthors,
			SubComments:   make([]entities.PostSubCommentItem, 0),
			CreatedAt:     createdAt,
		}

		filter = mongo_service.Filter().Eq("post_code", model.PostCode).Get()
		pushName = "comments"
	}

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$push": bson.M{pushName: entity},
		"$inc":  bson.M{"comments_count": 1},
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(model *models.PostCommentDelete) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	var filter bson.M
	var commentsField string

	if model.SubComment {
		filter = mongo_service.Filter().
			Eq("comments.comment_code", model.ParentCommentCode).
			Eq("comments.subcomments.comment_code", model.CommentCode).Get()

		commentsField = "comments.$.subcomments"
	} else {
		filter = mongo_service.Filter().Eq("comments.comment_code", model.CommentCode).Get()
		commentsField = "comments"
	}

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.M{
		"$pull": bson.M{
			commentsField: bson.M{
				"comment_code": model.CommentCode,
			},
		},
		"$inc": bson.M{"comments_count": -1},
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
