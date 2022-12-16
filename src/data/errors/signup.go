package errors

import "tok-core/src/data/models"

var (
	SignupBind = &models.Error{
		TechMessage:     "Signup model bind error",
		BusinessMessage: "Возника ошибка при регистрации",
	}

	SignupValidate = &models.Error{
		TechMessage:     "Signup model validate error",
		BusinessMessage: "Заполнены неверные данные",
	}

	Signup = &models.Error{
		TechMessage:     "Sign up error",
		BusinessMessage: "Возникла ошибка при регистрации",
	}
)
