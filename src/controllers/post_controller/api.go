package post_controller

import (
	"github.com/google/uuid"
	"github.com/lowl11/lazy-collection/array"
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazylog/layers"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

/*
	_add создание нового поста
*/
func (controller *Controller) _add(session *entities.ClientSession, model *models.PostAdd) *models.Error {
	logger := definition.Logger

	// расширенная модель
	extendedModel := &models.PostAddExtended{
		Base: model,
	}

	// создание кода поста
	postCode := uuid.New().String()

	// загружаем изображение поста если оно есть
	var picturePath *string

	if model.Picture != nil {
		if imageConfig, err := controller.image.UploadPostPicture(&models.PostPicture{
			Name:   model.Picture.Name,
			Buffer: model.Picture.Buffer,
		}, session.Username, postCode); err != nil {
			return errors.PostCreateUploadPicture.With(err)
		} else {
			picturePath = &imageConfig.Path
			extendedModel.ImageConfig = imageConfig
		}
	} else {
		extendedModel.ImageConfig = &models.ImageConfig{
			Width:  0,
			Height: 0,
		}
	}

	// заданная категория кастомная
	var customCategory *string
	if model.CustomCategory != nil {
		if categoryCode, err := controller.postCategoryRepo.Create(*extendedModel.Base.CustomCategory); err != nil {
			return errors.PostCategoryCreate.With(err)
		} else {
			customCategory = &categoryCode
		}
	}

	// создание поста
	if err := controller.postRepo.Create(extendedModel, session.Username, postCode, picturePath, customCategory); err != nil {
		logger.Error(err, "Create post error")
		return errors.PostCreate.With(err)
	}

	// создание поста в рекомендациях
	go func() {
		// построить путь к изображению (если есть) для эластика
		var elasticPicture *models.PostGetPicture
		if extendedModel.ImageConfig != nil && extendedModel.ImageConfig.Path != "" {
			postPicturePath := "/images/post/" + session.Username + "/post_" + postCode + "/" + extendedModel.ImageConfig.Path

			elasticPicture = &models.PostGetPicture{
				Path:   &postPicturePath,
				Width:  extendedModel.ImageConfig.Width,
				Height: extendedModel.ImageConfig.Height,
			}
		}

		// получаем имя категории
		var categoryName string
		if model.CustomCategory != nil {
			categoryName = *model.CustomCategory
		} else {
			category, _ := controller.postCategoryRepo.GetByCode(model.CategoryCode)
			if category != nil {
				categoryName = category.Name
			} else {
				// если вдруг не получилось найти категорию, вставляем сам код
				categoryName = model.CategoryCode
			}
		}

		// завести пост в рекоммендациях (через elastic)
		if err := controller.feed.AddExplore(&models.PostElasticAdd{
			Code:         postCode,
			Text:         extendedModel.Base.Text,
			Category:     extendedModel.Base.CategoryCode,
			CategoryName: categoryName,
			Picture:      elasticPicture,
			Author:       session.Username,
			CreatedAt:    time.Now(),

			Keys: []string{categoryName},
		}); err != nil {
			logger.Error(err, "Add post to explore error", layers.Elastic)
			return
		}
	}()

	return nil
}

/*
	_categories возвращает список всех категорий
*/
func (controller *Controller) _categories() ([]models.PostCategoryGet, *models.Error) {
	logger := definition.Logger

	// получение списка всех категорий
	categories, err := controller.postCategoryRepo.GetAll()
	if err != nil {
		logger.Error(err, "Get all post categories error")
		return nil, errors.PostCategoryGetList.With(err)
	}

	// обработка списка категорий для клиента
	list := make([]models.PostCategoryGet, 0, len(categories))
	for _, item := range categories {
		list = append(list, models.PostCategoryGet{
			Code: item.Code,
			Name: item.Name,
		})
	}

	return list, nil
}

/*
	_delete удаление поста с БД и эластика
*/
func (controller *Controller) _delete(code string) *models.Error {
	logger := definition.Logger

	// получаем пост по коду
	post, err := controller.postRepo.GetByCode(code)
	if err != nil {
		logger.Error(err, "Get post by code for deletion error", layers.Database)
		return errors.PostsGetByCode.With(err)
	}

	// ошибка если пост не найден
	if post == nil {
		return errors.PostNotFound
	}

	// удаление поста по коду в БД
	if err = controller.postRepo.DeleteByCode(code); err != nil {
		logger.Error(err, "Delete post by code error", layers.Database)
		return errors.PostDelete.With(err)
	}

	// удаление поста по коду в эластике
	if err = controller.feed.DeleteExplore(code); err != nil {
		logger.Error(err, "Delete post by code error", layers.Elastic)
		return errors.PostDelete.With(err)
	}

	return nil
}

/*
	_fillExplore заполнение данных для "рекомендаций"
*/
func (controller *Controller) _fillExplore() *models.Error {
	logger := definition.Logger

	posts, err := controller.postRepo.GetExplore()
	if err != nil {
		logger.Error(err, "Get explore posts error", layers.Database)
		return errors.PostsGetExplore.With(err)
	}

	if err = controller.feed.FillExplore(posts); err != nil {
		logger.Error(err, "Fill posts for explore error", layers.Elastic)
		return errors.ExploreFill.With(err)
	}

	return nil
}

/*
	_like поставить лайк для поста
*/
func (controller *Controller) _like(session *entities.ClientSession, model *models.PostLike) *models.Error {
	logger := definition.Logger

	// проверяем существует ли запись по лайкам
	like, err := controller.postLikeRepo.Get(model.PostCode)
	if err != nil {
		return errors.PostLikeGet.With(err)
	}

	if like == nil {
		if err = controller.postLikeRepo.Create(model, session.Username); err != nil {
			logger.Error(err, "Like post error", layers.Mongo)
			return errors.PostLike.With(err)
		}
	} else {
		if err = controller.postLikeRepo.Like(model.PostCode, session.Username); err != nil {
			logger.Error(err, "Like post error", layers.Mongo)
			return errors.PostLike.With(err)
		}
	}

	return nil
}

/*
	_unlike убрать лайк с поста
*/
func (controller *Controller) _unlike(session *entities.ClientSession, model *models.PostUnlike) *models.Error {
	logger := definition.Logger

	if err := controller.postLikeRepo.Unlike(model.PostCode, session.Username); err != nil {
		logger.Error(err, "Unlike post error", layers.Mongo)
		return errors.PostUnlike.With(err)
	}

	return nil
}

/*
	_getLikes получить список лайков (с авторами) поста
*/
func (controller *Controller) _getLikes(session *entities.ClientSession, postCode string) (*models.PostLikeGet, *models.Error) {
	logger := definition.Logger

	likes, err := controller.postLikeRepo.Get(postCode)
	if err != nil {
		logger.Error(err, "Get post likes error", layers.Mongo)
		return nil, errors.PostLikeGet.With(err)
	}

	if likes == nil {
		return nil, errors.PostNotFound
	}

	authorsList := array.NewWithList[string](likes.LikeAuthors...)

	return &models.PostLikeGet{
		LikesCount:  likes.LikesCount,
		LikeAuthors: likes.LikeAuthors,
		Liked:       authorsList.Contains(session.Username),
	}, nil
}

/*
	_getComment
*/
func (controller *Controller) _getComment(postCode string) (*models.PostCommentGet, *models.Error) {
	logger := definition.Logger

	postComments, err := controller.postCommentRepo.GetByPost(postCode)
	if err != nil {
		logger.Error(err, "Get comments list error", layers.Mongo)
		return nil, errors.PostCommentGet.With(err)
	}

	if postComments == nil {
		return nil, errors.PostCommentNotFound
	}

	commentsList := type_list.NewWithList[entities.PostCommentItem, models.PostCommentItem](postComments.Comments...)

	item := &models.PostCommentGet{
		PostCode:   postComments.PostCode,
		PostAuthor: postComments.PostAuthor,

		Comments: commentsList.Select(func(item entities.PostCommentItem) models.PostCommentItem {
			subCommentsList := type_list.NewWithList[entities.PostSubCommentItem, models.PostSubCommentItem](item.SubComments...)

			return models.PostCommentItem{
				CommentCode:   item.CommentCode,
				CommentAuthor: item.CommentAuthor,
				CommentText:   item.CommentText,
				LikesCount:    item.LikesCount,
				LikeAuthors:   item.LikeAuthors,
				CreatedAt:     item.CreatedAt,

				SubComments: subCommentsList.Select(func(item entities.PostSubCommentItem) models.PostSubCommentItem {
					return models.PostSubCommentItem{
						CommentCode:   item.CommentCode,
						CommentAuthor: item.CommentAuthor,
						CommentText:   item.CommentText,
						LikesCount:    item.LikesCount,
						LikeAuthors:   item.LikeAuthors,
					}
				}).Slice(),
			}
		}).Slice(),
	}

	return item, nil
}

/*
	_addComment
*/
func (controller *Controller) _addComment(session *entities.ClientSession, model *models.PostCommentAdd) (string, *models.Error) {
	logger := definition.Logger

	commentCode := uuid.New().String()

	// комментарий у поста может быть первым и не первым
	// если комментарий первый то значит записи в Mongo у поста нет (для комментариев)
	// значит, нужно создать его чтобы в дальнейшем в него писать
	if model.FirstComment {
		if err := controller.postCommentRepo.Create(model, session.Username, commentCode); err != nil {
			logger.Error(err, "Create new post comment error", layers.Mongo)
			return "", errors.PostCommentCreate.With(err)
		}
	} else {
		if err := controller.postCommentRepo.Append(model, session.Username, commentCode); err != nil {
			logger.Error(err, "Create new post comment error", layers.Mongo)
			return "", errors.PostCommentCreate.With(err)
		}
	}

	return commentCode, nil
}
