package errors

import "tok-core/src/data/models"

var (
	PostNotFound = &models.Error{
		TechMessage:     "Post not found",
		BusinessMessage: "Публикация не найдена",
	}

	PostsGetByUsernameParam = &models.Error{
		TechMessage:     "Username required",
		BusinessMessage: "Произошла ошибка",
	}
	PostsGetByCategoryParam = &models.Error{
		TechMessage:     "Category code required",
		BusinessMessage: "Произошла ошибка",
	}
	PostsGetBySingleParam = &models.Error{
		TechMessage:     "Post code required",
		BusinessMessage: "Произошла ошибка",
	}

	PostCreateBind = &models.Error{
		TechMessage:     "Create post model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	PostCreateValidate = &models.Error{
		TechMessage:     "Create post model validate error",
		BusinessMessage: "Произошла ошибка",
	}

	PostCreate = &models.Error{
		TechMessage:     "Create post error",
		BusinessMessage: "Произошла ошибка",
	}
	PostsGetByUsername = &models.Error{
		TechMessage:     "Get posts list by username error",
		BusinessMessage: "Произошла ошибка",
	}
	PostsGetByCategory = &models.Error{
		TechMessage:     "Get posts list by category error",
		BusinessMessage: "Произошла ошибка",
	}
	PostsGetByCode = &models.Error{
		TechMessage:     "Get post by code error",
		BusinessMessage: "Произошла ошибка",
	}
)
