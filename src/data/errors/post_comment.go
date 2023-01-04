package errors

import "tok-core/src/data/models"

var (
	PostCommentNotFound = &models.Error{
		TechMessage:     "Comments not found",
		BusinessMessage: "Комментарии не найдены",
	}
	PostCommentNotYours = &models.Error{
		TechMessage:     "This post is not yours",
		BusinessMessage: "Нельзя удалять чужие комментарии",
	}

	PostCommentCreateBind = &models.Error{
		TechMessage:     "Create post comment model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentDeleteBind = &models.Error{
		TechMessage:     "Delete post comment model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentLikeBind = &models.Error{
		TechMessage:     "Like post comment model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentUnlikeBind = &models.Error{
		TechMessage:     "Unlike post comment model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	PostCommentCreateValidation = &models.Error{
		TechMessage:     "Create post comment model validation error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentDeleteValidation = &models.Error{
		TechMessage:     "Delete post comment model validation error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentLikeValidation = &models.Error{
		TechMessage:     "Like post comment model validation error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentUnlikeValidation = &models.Error{
		TechMessage:     "Unlike post comment model validation error",
		BusinessMessage: "Произошла ошибка",
	}

	PostCommentGet = &models.Error{
		TechMessage:     "Get post comments list error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentCreate = &models.Error{
		TechMessage:     "Create new post comment error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentDelete = &models.Error{
		TechMessage:     "Delete post comment error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentLike = &models.Error{
		TechMessage:     "Like post comment error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentUnlike = &models.Error{
		TechMessage:     "Unlike post comment error",
		BusinessMessage: "Произошла ошибка",
	}
)
