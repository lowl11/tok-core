package errors

import "tok-core/src/data/models"

var (
	ProfileUpdateBind = &models.Error{
		TechMessage:     "Bind profile update error",
		BusinessMessage: "Возникла ошибка при изменении",
	}
	ProfileAvatarBind = &models.Error{
		TechMessage:     "Bind model upload avatar error",
		BusinessMessage: "Ошибка загрузки изображения",
	}
	ProfileWallpaperBind = &models.Error{
		TechMessage:     "Bind model upload wallpaper error",
		BusinessMessage: "Ошибка загрузки изображения",
	}

	ProfileUpdateValidate = &models.Error{
		TechMessage:     "Validate profile update error",
		BusinessMessage: "Введены неверные данные",
	}
	ProfileAvatarValidate = &models.Error{
		TechMessage:     "Upload avatar validation error",
		BusinessMessage: "Ошибка загрузки изображения",
	}
	ProfileWallpaperValidate = &models.Error{
		TechMessage:     "Upload wallpaper validation error",
		BusinessMessage: "Ошибка загрузки изображения",
	}

	ProfileUpdate = &models.Error{
		TechMessage:     "Update profile info error",
		BusinessMessage: "Возникла ошибка при изменении",
	}
	ProfileAvatar = &models.Error{
		TechMessage:     "Upload profile avatar error",
		BusinessMessage: "Ошибка загрузки изображения",
	}
	ProfileWallpaper = &models.Error{
		TechMessage:     "Upload profile wallpaper error",
		BusinessMessage: "Ошибка загрузки изображения",
	}
)
