package errors

import "tok-core/src/data/models"

var (
	ExploreFill = &models.Error{
		TechMessage:     "Fill explore feed error",
		BusinessMessage: "Произошла ошибка",
	}
)
