package errors

import "tok-core/src/data/models"

var (
	NotificationGetInfoParam = &models.Error{
		TechMessage:     "Username required",
		BusinessMessage: "Произошла ошибка",
	}
	NotificationGetCountParam = &models.Error{
		TechMessage:     "Username required",
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
)
