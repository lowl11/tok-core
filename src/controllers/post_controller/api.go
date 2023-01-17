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
	"tok-core/src/repositories/feed_repository"
	"tok-core/src/services/action_helper"
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
			logger.Error(err, "Create post category error", layers.Database)
			return errors.PostCategoryCreate.With(err)
		} else {
			customCategory = &categoryCode
			controller.category.Add(categoryCode, *extendedModel.Base.CustomCategory)
		}
	}

	// создание поста
	if err := controller.postRepo.Create(extendedModel, session.Username, postCode, picturePath, customCategory); err != nil {
		logger.Error(err, "Create post error", layers.Database)
		return errors.PostCreate.With(err)
	}

	// увеличиваем кол-во постов в категории
	go func() {
		if err := controller.categoryCountRepo.IncrementExist(model.CategoryCode); err != nil {
			logger.Error(err, "Increment categories count error", layers.Mongo)
		}
	}()

	// создание поста в рекомендациях
	go func() {
		var feedPostPicture *entities.FeedPostPicture

		if extendedModel.ImageConfig != nil {
			feedPostPicture = &entities.FeedPostPicture{
				Path:   extendedModel.ImageConfig.Path,
				Width:  extendedModel.ImageConfig.Width,
				Height: extendedModel.ImageConfig.Height,
			}
		}

		feedPost := &entities.FeedPost{
			PostCode:     postCode,
			PostCategory: model.CategoryCode,

			AuthorUsername: session.Username,
			AuthorName:     session.Name,
			AuthorAvatar:   session.Avatar,

			PostText: model.Text,
			Picture:  feedPostPicture,

			CreatedAt: time.Now(),
		}

		if err := controller.feedRepo.AddPostExist(feed_repository.Explore, feedPost); err != nil {
			logger.Error(err, "Add post to explore feed error", layers.Mongo)
		}

		if err := controller.feedRepo.AddPostExist(feed_repository.Unauthorized, feedPost); err != nil {
			logger.Error(err, "Add post to unauthorized feed error", layers.Mongo)
		}
	}()

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

	// удаление поста по коду в рекомендациях
	go func() {
		if err = controller.feedRepo.RemovePost(feed_repository.Explore, code); err != nil {
			logger.Error(err, "Remove post to explore feed error", layers.Mongo)
		}

		if err = controller.feedRepo.RemovePost(feed_repository.Unauthorized, code); err != nil {
			logger.Error(err, "Remove post to unauthorized feed error", layers.Mongo)
		}
	}()

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

	go func() {
		// отправляем уведомление
		if err = controller.notification.Push(action_helper.PostLike, model.PostAuthor, session.Username, &entities.NotificationBody{
			PostCode: &model.PostCode,
		}); err != nil {
			logger.Error(err, "Sending notification error", layers.Rabbit)
		}

		// запись интереса пользователя
		if err = controller.userInterestRepo.IncreaseCategoryExist(session.Username, model.PostCategory); err != nil {
			logger.Error(err, "Increase user category interest error", layers.Mongo)
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
		if err = controller.userInterestRepo.DecreaseCategoryExist(session.Username, model.PostCategory); err != nil {
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
			var mySubscription bool
			if array.NewWithList[string](session.Subscriptions.Subscriptions...).Contains(item.Username) {
				mySubscription = true
			}

			return models.UserDynamicGet{
				Username:   item.Username,
				Avatar:     item.Avatar,
				Name:       item.Name,
				Subscribed: &mySubscription,
			}
		}).Sub(from, from+10).Slice(),
		Liked: dynamicUsersList.Single(func(item entities.UserDynamicGet) bool {
			return item.Username == session.Username
		}) != nil,
	}, nil
}

/*
	_getComment получение комментариев под постом
*/
func (controller *Controller) _getComment(session *entities.ClientSession, postCode string, page int) (*models.PostCommentGet, *models.Error) {
	logger := definition.Logger

	from := (page - 1) * 15
	to := from + 15

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

	usernames := make([]string, 0, commentsList.Size())
	commentsList.Each(func(item entities.PostCommentItem) {
		usernames = append(usernames, item.CommentAuthor)
	})

	usersDynamic, err := controller.userRepo.GetDynamicByUsernames(usernames)
	if err != nil {
		return nil, errors.UserDynamicGet.With(err)
	}

	userList := type_list.NewWithList[entities.UserDynamicGet, models.UserDynamicGet](usersDynamic...).Select(func(item entities.UserDynamicGet) models.UserDynamicGet {
		return models.UserDynamicGet{
			Username: item.Username,
			Name:     item.Name,
			Avatar:   item.Avatar,
		}
	})

	item := &models.PostCommentGet{
		PostCode:   postComments.PostCode,
		PostAuthor: postComments.PostAuthor,

		Comments: commentsList.Select(func(item entities.PostCommentItem) models.PostCommentItem {
			subCommentsList := type_list.NewWithList[entities.PostSubCommentItem, models.PostSubCommentItem](item.SubComments...)

			return models.PostCommentItem{
				CommentCode:   item.CommentCode,
				CommentAuthor: userList.Single(func(user models.UserDynamicGet) bool { return user.Username == item.CommentAuthor }),
				CommentText:   item.CommentText,
				LikesCount:    item.LikesCount,
				MyLike:        array.NewWithList[string](item.LikeAuthors...).Contains(session.Username),
				CreatedAt:     item.CreatedAt,

				SubCommentsSize: subCommentsList.Size(),
			}
		}).Sub(from, to).Slice(),
	}

	return item, nil
}

/*
	_getSubcomment получение подкомментариев под родительским комментарием
*/
func (controller *Controller) _getSubcomment(session *entities.ClientSession, postCode, commentCode string, page int) (*models.PostSubcommentGet, *models.Error) {
	logger := definition.Logger

	from := (page - 1) * 15
	to := from + 15

	postComments, err := controller.postCommentRepo.GetByPost(postCode)
	if err != nil {
		logger.Error(err, "Get comments list error", layers.Mongo)
		return nil, errors.PostCommentGet.With(err)
	}

	if postComments == nil {
		return &models.PostSubcommentGet{
			PostCode:    postCode,
			PostAuthor:  "",
			Subcomments: make([]models.PostSubCommentItem, 0),
		}, nil
	}

	commentsList := type_list.NewWithList[entities.PostCommentItem, models.PostCommentItem](postComments.Comments...)
	subcomments := commentsList.Single(func(item entities.PostCommentItem) bool {
		return item.CommentCode == commentCode
	}).SubComments

	if subcomments == nil {
		return &models.PostSubcommentGet{
			PostCode:    postCode,
			PostAuthor:  "",
			Subcomments: make([]models.PostSubCommentItem, 0),
		}, nil
	}

	subCommentsList := type_list.NewWithList[entities.PostSubCommentItem, models.PostSubCommentItem](subcomments...)

	usernames := make([]string, 0, commentsList.Size())
	commentsList.Each(func(item entities.PostCommentItem) {
		usernames = append(usernames, item.CommentAuthor)
	})

	usersDynamic, err := controller.userRepo.GetDynamicByUsernames(usernames)
	if err != nil {
		return nil, errors.UserDynamicGet.With(err)
	}

	userList := type_list.NewWithList[entities.UserDynamicGet, models.UserDynamicGet](usersDynamic...).Select(func(item entities.UserDynamicGet) models.UserDynamicGet {
		return models.UserDynamicGet{
			Username: item.Username,
			Name:     item.Name,
			Avatar:   item.Avatar,
		}
	})

	item := &models.PostSubcommentGet{
		PostCode:   postComments.PostCode,
		PostAuthor: postComments.PostAuthor,

		Subcomments: subCommentsList.Select(func(item entities.PostSubCommentItem) models.PostSubCommentItem {
			return models.PostSubCommentItem{
				CommentCode:   item.CommentCode,
				CommentAuthor: userList.Single(func(user models.UserDynamicGet) bool { return user.Username == item.CommentAuthor }),
				CommentText:   item.CommentText,
				LikesCount:    item.LikesCount,
				MyLike:        array.NewWithList[string](item.LikeAuthors...).Contains(session.Username),
			}
		}).Sub(from, to).Slice(),
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
		1. Проверить, существует ли запись в feeds
		2. Если его нет, получить посты с БД за последний месяц и записать
	*/

	feed, err := controller.feedRepo.Get(feed_repository.Explore)
	if err != nil {
		return err
	}

	// если фида есть и он не пустой, ничего не делать
	if feed != nil && feed.Count > 0 {
		return nil
	}

	// если фида нет или он пустой, вытащить посты с БД и добавить в коллекцию
	list, err := controller.postRepo.GetExplore()
	if err != nil {
		return err
	}

	if err = controller.feedRepo.AddPostListExist(feed_repository.Explore, type_list.NewWithList[entities.PostGet, entities.FeedPost](list...).Select(
		func(item entities.PostGet) entities.FeedPost {
			var picture *entities.FeedPostPicture

			if item.Picture != nil {
				picture = &entities.FeedPostPicture{
					Path:   *item.Picture,
					Width:  item.PictureWidth,
					Height: item.PictureHeight,
				}
			}

			return entities.FeedPost{
				PostCode:     item.Code,
				PostCategory: item.CategoryCode,

				AuthorUsername: item.AuthorUsername,
				AuthorName:     item.AuthorName,
				AuthorAvatar:   item.AuthorAvatar,

				PostText: item.Text,
				Picture:  picture,

				CreatedAt: item.CreatedAt,
			}
		}).Slice()); err != nil {
		return err
	}

	return nil
}

/*
	_fillUnauthorizedFeed создание и заполнение индекса для главной ленты "неавторизованных"
*/
func (controller *Controller) _fillUnauthorizedFeed() error {
	/*
		1. Проверить, существует ли запись в feeds
		2. Если его нет, получить посты с БД за последний месяц и записать
	*/

	feed, err := controller.feedRepo.Get(feed_repository.Unauthorized)
	if err != nil {
		return err
	}

	// если фида есть и он не пустой, ничего не делать
	if feed != nil && feed.Count > 0 {
		return nil
	}

	// если фида нет или он пустой, вытащить посты с БД и добавить в коллекцию
	list, err := controller.postRepo.GetUnauthorized()
	if err != nil {
		return err
	}

	if err = controller.feedRepo.AddPostListExist(feed_repository.Unauthorized, type_list.NewWithList[entities.PostGet, entities.FeedPost](list...).Select(
		func(item entities.PostGet) entities.FeedPost {
			var picture *entities.FeedPostPicture

			if item.Picture != nil {
				picture = &entities.FeedPostPicture{
					Path:   *item.Picture,
					Width:  item.PictureWidth,
					Height: item.PictureHeight,
				}
			}

			return entities.FeedPost{
				PostCode:     item.Code,
				PostCategory: item.CategoryCode,

				AuthorUsername: item.AuthorUsername,
				AuthorName:     item.AuthorName,
				AuthorAvatar:   item.AuthorAvatar,

				PostText: item.Text,
				Picture:  picture,

				CreatedAt: item.CreatedAt,
			}
		}).Slice()); err != nil {
		return err
	}

	return nil
}
