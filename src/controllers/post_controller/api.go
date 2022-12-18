package post_controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) Add(ctx echo.Context) error {
	logger := definition.Logger
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateBind.With(err))
	}

	if err := controller.validatePostCreate(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateValidate.With(err))
	}

	postCode := uuid.New().String()
	if err := controller.postRepo.Create(&model, session.Username, postCode); err != nil {
		logger.Error(err, "Create post error")
		return controller.Error(ctx, errors.PostCreate.With(err))
	}

	return controller.Ok(ctx, postCode)
}

func (controller *Controller) Categories(ctx echo.Context) error {
	logger := definition.Logger

	categories, err := controller.postCategoryRepo.GetAll()
	if err != nil {
		logger.Error(err, "Get all post categories error")
		return controller.Error(ctx, errors.PostCategoryGetList.With(err))
	}

	list := make([]models.PostCategoryGet, 0, len(categories))
	for _, item := range categories {
		list = append(list, models.PostCategoryGet{
			Code: item.Code,
			Name: item.Name,
		})
	}

	return controller.Ok(ctx, list)
}
