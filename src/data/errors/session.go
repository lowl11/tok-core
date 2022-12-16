package errors

import "tok-core/src/data/models"

var (
	SessionNotFound = &models.Error{
		TechMessage:     "User session not found",
		BusinessMessage: "Возника ошибка при авторизации",
	}

	SessionCreate = &models.Error{
		TechMessage:     "Create user session error",
		BusinessMessage: "Возника техническая ошибка",
	}
	SessionUpdate = &models.Error{
		TechMessage:     "Update user session error",
		BusinessMessage: "Возника техническая ошибка",
	}

	SessionGet = &models.Error{
		TechMessage:     "Get user session error",
		BusinessMessage: "Возника ошибка при авторизации",
	}
)
