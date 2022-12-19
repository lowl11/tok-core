package errors

import "tok-core/src/data/models"

var (
	SearchUserBind = &models.Error{
		TechMessage:     "Search users model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	SearchUserValidate = &models.Error{
		TechMessage:     "Search users model validate error",
		BusinessMessage: "Произошла ошибка",
	}

	SearchUser = &models.Error{
		TechMessage:     "Search users error",
		BusinessMessage: "Произошла ошибка",
	}
)
