package post_comment_repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
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
				Text:          model.CommentText,
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
