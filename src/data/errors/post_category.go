package errors

import "tok-core/src/data/models"

var (
	PostCategoryNotFound = &models.Error{
		TechMessage:     "Post category not found",
		BusinessMessage: "Категория публикации не найдена",
	}

	PostCategoryCreate = &models.Error{
		TechMessage:     "Create post category error",
		BusinessMessage: "Произошла ошибка при создании категории",
	}
	PostCategoryGetList = &models.Error{
		TechMessage:     "Post categories error",
		BusinessMessage: "Произошла ошибка",
	}
)
