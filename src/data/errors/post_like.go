package errors

import "tok-core/src/data/models"

var (
	PostLike = &models.Error{
		TechMessage:     "Like post error",
		BusinessMessage: "Произошла ошибка",
	}
)
