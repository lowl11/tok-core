package errors

import "tok-core/src/data/models"

var (
	UserNotFound = &models.Error{
		TechMessage:     "User not found",
		BusinessMessage: "Введенные данные неверны",
	}

	UserGet = &models.Error{
		TechMessage:     "Get user error",
		BusinessMessage: "Ошибка получения пользователя",
	}
	UserSearch = &models.Error{
		TechMessage:     "Search users by query error",
		BusinessMessage: "Ошибка получения пользователя",
	}
	UserEncryptPassword = &models.Error{
		TechMessage:     "Encrypting password error",
		BusinessMessage: "Возника техническая ошибка",
	}
	UserDecryptPassword = &models.Error{
		TechMessage:     "Decrypting password error",
		BusinessMessage: "Возника техническая ошибка",
	}
)
