package post_comment_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
	"tok-core/src/services/mongo_service"
)

func (repo *Repository) GetAll() ([]entities.PostCommentGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	cursor, err := repo.connection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer repo.CloseCursor(cursor)

	list := make([]entities.PostCommentGet, 0)
	for cursor.Next(ctx) {
		item := entities.PostCommentGet{}
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	repo.LogError(cursor.Err())

	return list, nil
}

func (repo *Repository) Create(model *models.PostCommentAdd, commentCode string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	entity := entities.PostCommentCreate{
		PostCode:   model.PostCode,
		PostAuthor: model.PostAuthor,
		Comments: []entities.PostCommentItem{
			{
				CommentCode:   commentCode,
				CommentAuthor: model.CommentAuthor,
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

func (repo *Repository) Append(model *models.PostCommentAdd, commentCode string) error {
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
			CommentAuthor: model.CommentAuthor,
			CommentText:   model.CommentText,
			LikeAuthors:   likeAuthors,
			CreatedAt:     createdAt,
		}

		filter = mongo_service.Filter().Eq("post_code", model.PostCode).Eq("comments.comment_code", model.ParentCommentCode).Get()
		pushName = "comments.$.subcomments"
	} else {
		entity = entities.PostCommentAppendComment{
			CommentCode:   commentCode,
			CommentAuthor: model.CommentAuthor,
			CommentText:   model.CommentText,
			LikeAuthors:   likeAuthors,
			SubComments:   make([]entities.PostSubCommentItem, 0),
			CreatedAt:     createdAt,
		}

		filter = mongo_service.Filter().Eq("post_code", model.PostCode).Get()
		pushName = "comments"
	}

	if _, err := repo.connection.UpdateOne(ctx, filter, bson.D{
		{"$push", bson.D{
			{Key: pushName, Value: entity},
		}},
	}); err != nil {
		return err
	}

	return nil
}
