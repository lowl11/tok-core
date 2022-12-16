package errors

import "tok-core/src/data/models"

var (
	LoginBind = &models.Error{
		TechMessage:     "Login model bind error",
		BusinessMessage: "Возника ошибка при авторизации",
	}

	LoginValidate = &models.Error{
		TechMessage:     "Login model validate error",
		BusinessMessage: "Заполнены неверные данные",
	}

	Login = &models.Error{
		TechMessage:     "Login error",
		BusinessMessage: "Возника ошибка при авторизации",
	}
	LoginPassword = &models.Error{
		TechMessage:     "Wrong password",
		BusinessMessage: "Неверный пароль",
	}
)
