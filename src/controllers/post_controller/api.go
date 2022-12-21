package post_controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

/*
	Add создание нового поста
*/
func (controller *Controller) Add(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	// связка модели
	model := models.PostAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateBind.With(err))
	}

	// валидация модели
	if err := controller.validatePostCreate(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateValidate.With(err))
	}

	// создание кода поста
	postCode := uuid.New().String()

	// загружаем изображение поста если оно есть
	var picturePath string
	if model.Picture != nil {
		if uploadedPicturePath, err := controller.image.UploadPostPicture(&models.PostPicture{
			Name:   model.Picture.Name,
			Buffer: model.Picture.Buffer,
		}, session.Username, postCode); err != nil {
			return controller.Error(ctx, errors.PostCreateUploadPicture.With(err))
		} else {
			picturePath = uploadedPicturePath
		}
	}

	// заданная категория кастомная
	var customCategory *string
	if model.CustomCategory != nil {
		if categoryCode, err := controller.postCategoryRepo.Create(*model.CustomCategory); err != nil {
			return controller.Error(ctx, errors.PostCategoryCreate.With(err))
		} else {
			customCategory = &categoryCode
		}
	}

	// создание поста
	if err := controller.postRepo.Create(&model, session.Username, postCode, picturePath, customCategory); err != nil {
		logger.Error(err, "Create post error")
		return controller.Error(ctx, errors.PostCreate.With(err))
	}

	return controller.Ok(ctx, postCode)
}

/*
	Categories возвращает список всех категорий
*/
func (controller *Controller) Categories(ctx echo.Context) error {
	logger := definition.Logger

	// получение списка всех категорий
	categories, err := controller.postCategoryRepo.GetAll()
	if err != nil {
		logger.Error(err, "Get all post categories error")
		return controller.Error(ctx, errors.PostCategoryGetList.With(err))
	}

	// обработка списка категорий для клиента
	list := make([]models.PostCategoryGet, 0, len(categories))
	for _, item := range categories {
		list = append(list, models.PostCategoryGet{
			Code: item.Code,
			Name: item.Name,
		})
	}

	return controller.Ok(ctx, list)
}
