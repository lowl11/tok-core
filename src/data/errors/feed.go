package errors

import "tok-core/src/data/models"

var (
	FeedNotFound = &models.Error{
		TechMessage:     "Feed not found",
		BusinessMessage: "Произошла ошибка",
	}

	FeedGet = &models.Error{
		TechMessage:     "Get feed error",
		BusinessMessage: "Произошла ошибка",
	}
)
