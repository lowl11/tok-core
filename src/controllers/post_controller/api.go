package post_controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lowl11/lazy-collection/array"
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

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
	_addCategory создать категорию, код категории генерируется сам
*/
func (controller *Controller) _addCategory(name string) (string, *models.Error) {
	logger := definition.Logger

	code, err := controller.postCategoryRepo.Create(name)
	if err != nil {
		logger.Error(err, "Create post category error", layers.Database)
		return "", errors.PostCategoryCreate.With(err)
	}

	return code, nil
}

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

	// увеличиваем кол-во постов в категории
	go func() {
		categoryCount, err := controller.categoryCountRepo.Get(model.CategoryCode)
		if err != nil {
			logger.Error(err, "Get category count error", layers.Mongo)
			return
		}

		if categoryCount == nil {
			if err = controller.categoryCountRepo.Create(model.CategoryCode); err != nil {
				logger.Error(err, "Create category count error", layers.Mongo)
				return
			}
		} else {
			if err = controller.categoryCountRepo.Increment(model.CategoryCode); err != nil {
				logger.Error(err, "Increment categories count error", layers.Mongo)
			}
		}
	}()

	// создание поста в рекомендациях
	// TODO: implement

	return nil
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
	// TODO: implement me

	// уменьшаем кол-во постов в категории
	go func() {
		categoryCount, err := controller.categoryCountRepo.Get(code)
		if err != nil {
			logger.Error(err, "Get category count error", layers.Mongo)
			return
		}

		if categoryCount != nil && categoryCount.Count > 0 {
			if err = controller.categoryCountRepo.Decrement(code); err != nil {
				logger.Error(err, "Decrement categories count error", layers.Mongo)
				return
			}
		}
	}()

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
		// если вдруг тот кто авторизовался уже поставил лайк
		if array.NewWithList[string](like.LikeAuthors...).Contains(session.Username) {
			return nil
		}

		if err = controller.postLikeRepo.Like(model.PostCode, session.Username); err != nil {
			logger.Error(err, "Like post error", layers.Mongo)
			return errors.PostLike.With(err)
		}
	}

	// запись интереса пользователя
	go func() {
		interest, err := controller.userInterest.Get(session.Username)
		if err != nil {
			logger.Error(err, "Get user interest error", layers.Mongo)
			return
		}

		if interest == nil {
			if err = controller.userInterest.Create(session.Username, model.PostCategory); err != nil {
				logger.Error(err, "Create user interest error", layers.Mongo)
			}
		} else {
			if err = controller.userInterest.IncreaseCategory(session.Username, model.PostCategory); err != nil {
				logger.Error(err, "Increase user category interest error", layers.Mongo)
			}
		}
	}()

	return nil
}

/*
	_unlike убрать лайк с поста
*/
func (controller *Controller) _unlike(session *entities.ClientSession, model *models.PostUnlike) *models.Error {
	logger := definition.Logger

	// проверяем существует ли запись по лайкам
	like, err := controller.postLikeRepo.Get(model.PostCode)
	if err != nil {
		return errors.PostLikeGet.With(err)
	}

	// проверить вдруг кол-во лайков уже 0
	if like != nil && like.LikesCount == 0 {
		return nil
	}

	if err = controller.postLikeRepo.Unlike(model.PostCode, session.Username); err != nil {
		logger.Error(err, "Unlike post error", layers.Mongo)
		return errors.PostUnlike.With(err)
	}

	// запись интереса пользователя
	go func() {
		interest, err := controller.userInterest.Get(session.Username)
		if err != nil {
			logger.Error(err, "Get user interest error", layers.Mongo)
			return
		}

		// если записи нет, то и уменьшать не нужно
		if interest == nil {
			return
		}

		interestCategory := array.NewWithList[entities.UserInterestCategory](interest.Categories...).Single(func(item entities.UserInterestCategory) bool {
			return item.CategoryCode == model.PostCategory
		})

		// если категории нет, то и уменьшать не нужно
		if interestCategory == nil || interestCategory.Interest <= 0 {
			fmt.Println("i am here")
			return
		}

		if err = controller.userInterest.DecreaseCategory(session.Username, model.PostCategory); err != nil {
			logger.Error(err, "Decrease user category interest error", layers.Mongo)
		}
	}()

	return nil
}

/*
	_getLikes получить список лайков (с авторами) поста
*/
func (controller *Controller) _getLikes(session *entities.ClientSession, postCode string, page int) (*models.PostLikeGet, *models.Error) {
	logger := definition.Logger

	likes, err := controller.postLikeRepo.Get(postCode)
	if err != nil {
		logger.Error(err, "Get post likes error", layers.Mongo)
		return nil, errors.PostLikeGet.With(err)
	}

	if likes == nil {
		return &models.PostLikeGet{
			LikesCount:  0,
			LikeAuthors: make([]models.UserDynamicGet, 0),
			Liked:       false,
		}, nil
	}

	dynamicUsers, err := controller.userRepo.GetDynamicByUsernames(likes.LikeAuthors)
	if err != nil {
		logger.Error(err, "Get users dynamic info error", layers.Database)
		return nil, errors.UserDynamicGet.With(err)
	}

	dynamicUsersList := type_list.NewWithList[entities.UserDynamicGet, models.UserDynamicGet](dynamicUsers...)

	// высчитаем с какого по какой лайк нужен
	from := (page - 1) * 10
	if from > dynamicUsersList.Size() {
		from = 0
	}

	return &models.PostLikeGet{
		LikesCount: likes.LikesCount,
		LikeAuthors: dynamicUsersList.Select(func(item entities.UserDynamicGet) models.UserDynamicGet {
			return models.UserDynamicGet{
				Username: item.Username,
				Avatar:   item.Avatar,
				Name:     item.Name,
			}
		}).Sub(from, from+10).Slice(),
		Liked: dynamicUsersList.Single(func(item entities.UserDynamicGet) bool {
			return item.Username == session.Username
		}) != nil,
	}, nil
}

/*
	_getComment получение всего дерева комментариев под постом
*/
func (controller *Controller) _getComment(postCode string) (*models.PostCommentGet, *models.Error) {
	logger := definition.Logger

	postComments, err := controller.postCommentRepo.GetByPost(postCode)
	if err != nil {
		logger.Error(err, "Get comments list error", layers.Mongo)
		return nil, errors.PostCommentGet.With(err)
	}

	if postComments == nil {
		return &models.PostCommentGet{
			PostCode:   postCode,
			PostAuthor: "",
			Comments:   make([]models.PostCommentItem, 0),
		}, nil
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
	_addComment добавление нового комментария под постом
	Добавляется, как обычный комментарий, так и "подкомментарий"
	Уровень подкомментариев может быть максимум 1, то есть кто-то написал первый комментарий под постом
	И тот кто ответит, ему оставит "подкомментарий", это и будет последним уровнем
*/
func (controller *Controller) _addComment(session *entities.ClientSession, model *models.PostCommentAdd) (string, *models.Error) {
	logger := definition.Logger

	commentCode := uuid.New().String()

	// если клиент говорит что это первый коммент под постом
	// нужно проверить существует ли под этот пост запись
	if model.FirstComment {
		comment, err := controller.postCommentRepo.GetByPost(model.PostCode)
		if err != nil {
			return "", errors.PostCommentGet.With(err)
		}

		if comment != nil {
			model.FirstComment = false
		}
	}

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

/*
	_deleteComment удаление комментария или подкомментария
*/
func (controller *Controller) _deleteComment(session *entities.ClientSession, model *models.PostCommentDelete) *models.Error {
	logger := definition.Logger

	post, err := controller.postCommentRepo.GetByCode(model.CommentCode, model.SubComment)
	if err != nil {
		logger.Error(err, "Get post comment error", layers.Mongo)
		return errors.PostCommentGet.With(err)
	}

	if post == nil {
		return errors.PostCommentNotFound
	}

	var authorMatch bool
	for _, comment := range post.Comments {
		if len(comment.SubComments) > 0 {
			// TODO: что делать если удаляется комментарий у которого есть подкомментарии?
		}

		if session.Username == comment.CommentAuthor {
			authorMatch = true
			break
		}

		for _, subComment := range comment.SubComments {
			if session.Username == subComment.CommentAuthor {
				authorMatch = true
				break
			}
		}
	}

	if !authorMatch {
		return errors.PostCommentNotYours
	}

	if err = controller.postCommentRepo.Delete(model); err != nil {
		logger.Error(err, "Delete post comment error", layers.Mongo)
		return errors.PostCommentDelete.With(err)
	}

	return nil
}

/*
	_likeComment поставить лайк комментария под постом
*/
func (controller *Controller) _likeComment(session *entities.ClientSession, model *models.PostCommentLike) *models.Error {
	logger := definition.Logger

	if err := controller.postCommentRepo.Like(model, session.Username); err != nil {
		logger.Error(err, "Like post comment error", layers.Mongo)
		return errors.PostCommentLike.With(err)
	}

	return nil
}

/*
	_unlikeComment убрать лайк комментария под постом
*/
func (controller *Controller) _unlikeComment(session *entities.ClientSession, model *models.PostCommentUnlike) *models.Error {
	logger := definition.Logger

	if err := controller.postCommentRepo.Unlike(model, session.Username); err != nil {
		logger.Error(err, "Unlike post comment error", layers.Mongo)
		return errors.PostCommentUnlike.With(err)
	}

	return nil
}

/*
	_fillExploreFeed создание и заполнение индекса для "рекомендаций"
*/
func (controller *Controller) _fillExploreFeed() error {
	/*
		1. Проверить, существует ли сегодняшний индекс
		2. Если его нет, получить посты с БД
		3. Создать завтрашний индекс получив сегодняшние посты
	*/

	return nil
}

/*
	_fillUnauthorizedFeed создание и заполнение индекса для главной ленты "неавторизованных"
*/
func (controller *Controller) _fillUnauthorizedFeed() error {
	/*
		1. Проверить, существует ли сегодняшний индекс
		2. Если его нет, ничего не делать
		3. Создать завтрашний индекс получив сегодняшние посты
	*/

	return nil
}
