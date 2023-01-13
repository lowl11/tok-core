package feed_controller

import (
	"github.com/lowl11/lazy-collection/array"
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
	"tok-core/src/repositories/feed_repository"
	"tok-core/src/services/feed_helper"
)

/*
	_general лента на главной странице для неавторизованных
*/
func (controller *Controller) _general(page int) ([]models.PostGet, *models.Error) {
	//logger := definition.Logger
	return nil, nil
}

/*
	_main лента на главной странице
*/
func (controller *Controller) _main(session *entities.ClientSession, page int) ([]models.PostGet, *models.Error) {
	logger := definition.Logger

	// список подписок из сессии
	subscriptions := session.Subscriptions.Subscriptions

	// добавить себя в список
	subscriptions = append(subscriptions, session.Username)

	// посты с массивом юзернеймов из подписок
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByUsernameList(subscriptions, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by username list error")
		return nil, errors.PostsGetByUsernameList.With(err)
	}

	// массив с кодами постов
	postCodes := type_list.NewWithList[entities.PostGet, string](posts...).Select(func(item entities.PostGet) string {
		return item.Code
	}).Slice()

	// получаем лайки и комментарии постов
	likeArray, commentArray := feed_helper.LikesAndComments(postCodes)

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		// лайки
		var myLike bool
		var likeCount int

		if likeArray != nil {
			foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundLike != nil {
				likeCount = foundLike.LikesCount

				// поставил ли авторизовавшийся лайк
				likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
				myLike = likeAuthors.Contains(session.Username)
			}
		}

		// комментарии
		var commentCount int

		if commentArray != nil {
			foundComment := commentArray.Single(func(iterator entities.PostCommentGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundComment != nil {
				commentCount = foundComment.CommentsCount
			}
		}

		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,

			LikeCount:    likeCount,
			MyLike:       myLike,
			CommentCount: commentCount,
		})
	}

	return list, nil
}

/*
	_explore лента "рекомендаций"
*/
func (controller *Controller) _explore(session *entities.ClientSession, page int) ([]models.PostGet, *models.Error) {
	logger := definition.Logger

	// пагинация (от и до)
	from := (page - 1) * 10
	to := from + 10

	feed, err := controller.feedRepo.Get(feed_repository.Explore)
	if err != nil {
		logger.Error(err, "Get feed error", layers.Mongo)
		return nil, errors.FeedGet.With(err)
	}

	if feed == nil {
		return nil, errors.FeedNotFound
	}

	posts := array.NewWithList[entities.FeedPost](feed.Posts...).Shuffle().Sub(from, to).Slice()

	// массив с кодами постов
	postCodes := type_list.NewWithList[entities.FeedPost, string](posts...).Select(func(item entities.FeedPost) string {
		return item.PostCode
	}).Slice()

	// получаем лайки и комментарии постов
	likeArray, commentArray := feed_helper.LikesAndComments(postCodes)

	return type_list.NewWithList[entities.FeedPost, models.PostGet](posts...).Select(func(item entities.FeedPost) models.PostGet {
		var picture *models.PostGetPicture

		if item.Picture != nil {
			picture = &models.PostGetPicture{
				Path:   &item.Picture.Path,
				Width:  item.Picture.Width,
				Height: item.Picture.Height,
			}
		}

		// лайки
		var myLike bool
		var likeCount int

		if likeArray != nil {
			foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
				return iterator.PostCode == item.PostCode
			})

			if foundLike != nil {
				likeCount = foundLike.LikesCount

				// поставил ли авторизовавшийся лайк
				likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
				myLike = likeAuthors.Contains(session.Username)
			}
		}

		// комментарии
		var commentCount int

		if commentArray != nil {
			foundComment := commentArray.Single(func(iterator entities.PostCommentGetList) bool {
				return iterator.PostCode == item.PostCode
			})

			if foundComment != nil {
				commentCount = foundComment.CommentsCount
			}
		}

		// категория поста
		var categoryName string
		category := controller.category.Get(item.PostCategory)
		if category != nil {
			categoryName = category.Name
		}

		return models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.PostCategory,
			CategoryName: categoryName,

			Code:      item.PostCode,
			Text:      item.PostText,
			Picture:   picture,
			CreatedAt: item.CreatedAt,

			LikeCount:    likeCount,
			MyLike:       myLike,
			CommentCount: commentCount,
		}
	}).Slice(), nil
}

/*
	_user лента на странице пользователя
*/
func (controller *Controller) _user(session *entities.ClientSession, username string, page int) ([]models.PostGet, *models.Error) {
	logger := definition.Logger

	// получить список постов по логину
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByUsername(username, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by username error")
		return nil, errors.PostsGetByUsername.With(err)
	}

	// массив с кодами постов
	postCodes := type_list.NewWithList[entities.PostGet, string](posts...).Select(func(item entities.PostGet) string {
		return item.Code
	}).Slice()

	// получаем лайки и комментарии постов
	likeArray, commentArray := feed_helper.LikesAndComments(postCodes)

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		var myLike bool
		var likeCount int

		if likeArray != nil {
			foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundLike != nil {
				likeCount = foundLike.LikesCount

				likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
				myLike = likeAuthors.Contains(session.Username)
			}
		}

		// комментарии
		var commentCount int

		if commentArray != nil {
			foundComment := commentArray.Single(func(iterator entities.PostCommentGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundComment != nil {
				commentCount = foundComment.CommentsCount
			}
		}

		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,

			LikeCount:    likeCount,
			MyLike:       myLike,
			CommentCount: commentCount,
		})
	}

	return list, nil
}

/*
	_category лента на странице категории
*/
func (controller *Controller) _category(session *entities.ClientSession, categoryCode string, page int) ([]models.PostGet, *models.Error) {
	logger := definition.Logger

	// получить список постов по логину
	offset := (page - 1) * 10
	size := 10
	posts, err := controller.postRepo.GetByCategory(categoryCode, offset, size)
	if err != nil {
		logger.Error(err, "Get posts list by category error")
		return nil, errors.PostsGetByCategory.With(err)
	}

	// массив с кодами постов
	postCodes := type_list.NewWithList[entities.PostGet, string](posts...).Select(func(item entities.PostGet) string {
		return item.Code
	}).Slice()

	// получаем лайки и комментарии постов
	likeArray, commentArray := feed_helper.LikesAndComments(postCodes)

	// обработка списка постов для клиента
	list := make([]models.PostGet, 0, len(posts))
	for _, item := range posts {
		var myLike bool
		var likeCount int

		if likeArray != nil {
			foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundLike != nil {
				likeCount = foundLike.LikesCount

				likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
				myLike = likeAuthors.Contains(session.Username)
			}
		}

		// комментарии
		var commentCount int

		if commentArray != nil {
			foundComment := commentArray.Single(func(iterator entities.PostCommentGetList) bool {
				return iterator.PostCode == item.Code
			})

			if foundComment != nil {
				commentCount = foundComment.CommentsCount
			}
		}

		list = append(list, models.PostGet{
			AuthorUsername: item.AuthorUsername,
			AuthorName:     item.AuthorName,
			AuthorAvatar:   item.AuthorAvatar,

			CategoryCode: item.CategoryCode,
			CategoryName: item.CategoryName,

			Code: item.Code,
			Text: item.Text,
			Picture: &models.PostGetPicture{
				Path:   item.Picture,
				Width:  item.PictureWidth,
				Height: item.PictureHeight,
			},
			CreatedAt: item.CreatedAt,

			LikeCount:    likeCount,
			MyLike:       myLike,
			CommentCount: commentCount,
		})
	}

	return list, nil
}

/*
	Single страница одного поста
*/
func (controller *Controller) _single(session *entities.ClientSession, code string) (*models.PostGet, *models.Error) {
	logger := definition.Logger

	// получение списка постов
	post, err := controller.postRepo.GetByCode(code)
	if err != nil {
		logger.Error(err, "Get post by code error")
		return nil, errors.PostsGetByCode.With(err)
	}

	// если пост не найден
	if post == nil {
		return nil, errors.PostNotFound
	}

	// получаем лайки и комментарии постов
	likeArray, commentArray := feed_helper.LikesAndComments([]string{post.Code})

	var myLike bool
	var likeCount int

	if likeArray != nil {
		foundLike := likeArray.Single(func(iterator entities.PostLikeGetList) bool {
			return iterator.PostCode == post.Code
		})

		if foundLike != nil {
			likeCount = foundLike.LikesCount

			likeAuthors := array.NewWithList[string](foundLike.LikeAuthors...)
			myLike = likeAuthors.Contains(session.Username)
		}
	}

	// комментарии
	var commentCount int

	if commentArray != nil {
		foundComment := commentArray.Single(func(iterator entities.PostCommentGetList) bool {
			return iterator.PostCode == post.Code
		})

		if foundComment != nil {
			commentCount = foundComment.CommentsCount
		}
	}

	// обработка поста для клиента
	return &models.PostGet{
		AuthorUsername: post.AuthorUsername,
		AuthorName:     post.AuthorName,
		AuthorAvatar:   post.AuthorAvatar,

		CategoryCode: post.CategoryCode,
		CategoryName: post.CategoryName,

		Code: post.Code,
		Text: post.Text,
		Picture: &models.PostGetPicture{
			Path:   post.Picture,
			Width:  post.PictureWidth,
			Height: post.PictureHeight,
		},
		CreatedAt: post.CreatedAt,

		LikeCount:    likeCount,
		MyLike:       myLike,
		CommentCount: commentCount,
	}, nil
}
