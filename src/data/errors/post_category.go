package errors

import "tok-core/src/data/models"

var (
	PostCategoryNotFound = &models.Error{
		TechMessage:     "Post category not found",
		BusinessMessage: "Категория публикации не найдена",
	}

	PostCategoryCreateBind = &models.Error{
		TechMessage:     "Create post category bind model error",
		BusinessMessage: "Произошла ошибка при создании категории",
	}

	PostCategoryCreateValidate = &models.Error{
		TechMessage:     "Create post category bind model error",
		BusinessMessage: "Произошла ошибка при создании категории",
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
