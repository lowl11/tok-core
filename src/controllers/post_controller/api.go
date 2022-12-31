package post_controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	AddREST обертка для _add
*/
func (controller *Controller) AddREST(ctx echo.Context) error {
	// связка модели
	model := models.PostAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateBind.With(err))
	}

	// валидация модели
	if err := controller.validatePostCreate(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateValidate.With(err))
	}

	session := ctx.Get("client_session").(*entities.ClientSession)

	if err := controller._add(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
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
	CategoriesREST обертка для _categories
*/
func (controller *Controller) CategoriesREST(ctx echo.Context) error {
	list, err := controller._categories()
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, list)
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
	DeleteREST обертка для _delete
*/
func (controller *Controller) DeleteREST(ctx echo.Context) error {
	code := ctx.Param("code")
	if code == "" {
		return controller.Error(ctx, errors.PostDeleteParam)
	}

	if err := controller._delete(code); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
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
	FillExploreREST обертка для _fillExplore
*/
func (controller *Controller) FillExploreREST(ctx echo.Context) error {
	if err := controller._fillExplore(); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

func (controller *Controller) _like(model *models.PostLike) *models.Error {
	//

	return nil
}

func (controller *Controller) _unlike(model *models.PostUnlike) *models.Error {
	//

	return nil
}

func (controller *Controller) LikeREST(ctx echo.Context) error {
	model := models.PostLike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeBind.With(err))
	}

	if err := controller.validateLike(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeValidate.With(err))
	}

	if err := controller._like(&model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

func (controller *Controller) UnlikeREST(ctx echo.Context) error {
	model := models.PostUnlike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeBind.With(err))
	}

	if err := controller.validateUnlike(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeValidate.With(err))
	}

	if err := controller._unlike(&model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

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

func (controller *Controller) GetCommentREST(ctx echo.Context) error {
	postCode := ctx.QueryParam("code")

	comments, err := controller._getComment(postCode)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, comments)
}

func (controller *Controller) _addComment(model *models.PostCommentAdd) (string, *models.Error) {
	logger := definition.Logger

	commentCode := uuid.New().String()

	// комментарий у поста может быть первым и не первым
	// если комментарий первый то значит записи в Mongo у поста нет (для комментариев)
	// значит, нужно создать его чтобы в дальнейшем в него писать
	if model.FirstComment {
		if err := controller.postCommentRepo.Create(model, commentCode); err != nil {
			logger.Error(err, "Create new post comment error", layers.Mongo)
			return "", errors.PostCommentCreate.With(err)
		}
	} else {
		if err := controller.postCommentRepo.Append(model, commentCode); err != nil {
			logger.Error(err, "Create new post comment error", layers.Mongo)
			return "", errors.PostCommentCreate.With(err)
		}
	}

	return commentCode, nil
}

func (controller *Controller) AddCommentREST(ctx echo.Context) error {
	model := models.PostCommentAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeBind.With(err))
	}

	if err := controller.validateAddComment(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeValidate.With(err))
	}

	commentCode, err := controller._addComment(&model)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, commentCode)
}
