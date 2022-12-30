package errors

import "tok-core/src/data/models"

var (
	PostCommentGet = &models.Error{
		TechMessage:     "Get post comments list error",
		BusinessMessage: "Произошла ошибка",
	}
	PostCommentCreate = &models.Error{
		TechMessage:     "Create new post comment error",
		BusinessMessage: "Произошла ошибка",
	}
)
