package errors

import "tok-core/src/data/models"

var (
	PostLikeNotFound = &models.Error{
		TechMessage:     "Post likes not found",
		BusinessMessage: "Произошла ошибка",
	}

	PostLikeGetParam = &models.Error{
		TechMessage:     "Post Code required",
		BusinessMessage: "Произошла ошибка",
	}

	PostLikeBind = &models.Error{
		TechMessage:     "Like post model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	PostUnlikeBind = &models.Error{
		TechMessage:     "Unlike  post model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	PostLikeValidate = &models.Error{
		TechMessage:     "Like post model validate error",
		BusinessMessage: "Произошла ошибка",
	}
	PostUnlikeValidate = &models.Error{
		TechMessage:     "Unlike post model validate error",
		BusinessMessage: "Произошла ошибка",
	}

	PostLike = &models.Error{
		TechMessage:     "Like post error",
		BusinessMessage: "Произошла ошибка",
	}
	PostUnlike = &models.Error{
		TechMessage:     "Unlike post error",
		BusinessMessage: "Произошла ошибка",
	}
	PostLikeGet = &models.Error{
		TechMessage:     "Get post likes error",
		BusinessMessage: "Произошла ошибка",
	}
	PostLikeDelete = &models.Error{
		TechMessage:     "Delete post like(s) error",
		BusinessMessage: "Произошла ошибка",
	}
)
