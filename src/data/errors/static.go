package errors

import "tok-core/src/data/models"

var (
	RouteNotFound = &models.Error{
		TechMessage:     "Route not found",
		BusinessMessage: "Путь не найден",
	}

	Timeout = &models.Error{
		TechMessage:     "Request reached timed out",
		BusinessMessage: "Время работы истекло",
	}
)