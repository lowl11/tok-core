package errors

import "tok-core/src/data/models"

var (
	SearchUserBind = &models.Error{
		TechMessage:     "Search users model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchCategoryBind = &models.Error{
		TechMessage:     "Search categories model bind error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchSmartBind = &models.Error{
		TechMessage:     "Search smart model bind error",
		BusinessMessage: "Произошла ошибка",
	}

	SearchUserValidate = &models.Error{
		TechMessage:     "Search users model validate error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchCategoryValidate = &models.Error{
		TechMessage:     "Search categories model validate error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchSmartValidate = &models.Error{
		TechMessage:     "Search smart model validate error",
		BusinessMessage: "Произошла ошибка",
	}

	SearchUser = &models.Error{
		TechMessage:     "Search users error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchCategory = &models.Error{
		TechMessage:     "Search categories error",
		BusinessMessage: "Произошла ошибка",
	}
	SearchSmart = &models.Error{
		TechMessage:     "Search smart error",
		BusinessMessage: "Произошла ошибка",
	}
)
