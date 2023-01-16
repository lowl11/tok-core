package errors

import "tok-core/src/data/models"

var (
	NotificationReadBind = &models.Error{
		TechMessage:     "Read notifications model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	NotificationReadValidation = &models.Error{
		TechMessage:     "Read notifications model validation error",
		BusinessMessage: "Произошла ошибка",
	}

	NotificationGetInfo = &models.Error{
		TechMessage:     "Get notifications info error",
		BusinessMessage: "Произошла ошибка",
	}
	NotificationGetCount = &models.Error{
		TechMessage:     "Get notifications count error",
		BusinessMessage: "Произошла ошибка",
	}
	NotificationRead = &models.Error{
		TechMessage:     "Read notifications error",
		BusinessMessage: "Произошла ошибка",
	}
)
