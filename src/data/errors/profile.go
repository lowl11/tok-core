package errors

import "tok-core/src/data/models"

var (
	ProfileUpdateBind = &models.Error{
		TechMessage:     "Bind profile update error",
		BusinessMessage: "Возникла ошибка при изменении",
	}

	ProfileUpdateValidate = &models.Error{
		TechMessage:     "Validate profile update error",
		BusinessMessage: "Введены неверные данные",
	}

	ProfileUpdate = &models.Error{
		TechMessage:     "Update profile info error",
		BusinessMessage: "Возникла ошибка при изменении",
	}
)
