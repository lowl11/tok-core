package errors

import "tok-core/src/data/models"

var (
	UserIpUserNotFound = &models.Error{
		TechMessage:     "User not found",
		BusinessMessage: "Произошла ошибка",
	}

	UserIpGetByIp = &models.Error{
		TechMessage:     "Get username by IP address error",
		BusinessMessage: "Произошла ошибка",
	}
)
