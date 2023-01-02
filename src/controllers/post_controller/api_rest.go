package post_controller

import (
	"github.com/labstack/echo/v4"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
)

/*
	AddREST обертка для _add
*/
func (controller *Controller) AddREST(ctx echo.Context) error {
	// связка модели
	model := models.PostAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateBind.With(err))
	}

	// валидация модели
	if err := controller.validatePostCreate(&model); err != nil {
		return controller.Error(ctx, errors.PostCreateValidate.With(err))
	}

	session := ctx.Get("client_session").(*entities.ClientSession)

	if err := controller._add(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	CategoriesREST обертка для _categories
*/
func (controller *Controller) CategoriesREST(ctx echo.Context) error {
	list, err := controller._categories()
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, list)
}

/*
	DeleteREST обертка для _delete
*/
func (controller *Controller) DeleteREST(ctx echo.Context) error {
	code := ctx.Param("code")
	if code == "" {
		return controller.Error(ctx, errors.PostDeleteParam)
	}

	if err := controller._delete(code); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	FillExploreREST обертка для _fillExplore
*/
func (controller *Controller) FillExploreREST(ctx echo.Context) error {
	if err := controller._fillExplore(); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	LikeREST обертка для _like
*/
func (controller *Controller) LikeREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostLike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeBind.With(err))
	}

	if err := controller.validateLike(&model); err != nil {
		return controller.Error(ctx, errors.PostLikeValidate.With(err))
	}

	if err := controller._like(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	UnlikeREST обертка для _unlike
*/
func (controller *Controller) UnlikeREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostUnlike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostUnlikeBind.With(err))
	}

	if err := controller.validateUnlike(&model); err != nil {
		return controller.Error(ctx, errors.PostUnlikeValidate.With(err))
	}

	if err := controller._unlike(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	GetLikesREST обертка для _getLikes
*/
func (controller *Controller) GetLikesREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	postCode := ctx.Param("code")
	if postCode == "" {
		return errors.PostLikeGetParam
	}

	likes, err := controller._getLikes(session, postCode)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, likes)
}

/*
	GetCommentREST обертка для _getComment
*/
func (controller *Controller) GetCommentREST(ctx echo.Context) error {
	postCode := ctx.QueryParam("code")

	comments, err := controller._getComment(postCode)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, comments)
}

/*
	AddCommentREST обертка для _addComment
*/
func (controller *Controller) AddCommentREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostCommentAdd{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentCreateBind.With(err))
	}

	if err := controller.validateAddComment(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentCreateValidation.With(err))
	}

	commentCode, err := controller._addComment(session, &model)
	if err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, commentCode)
}

/*
	DeleteCommentREST обертка для _deleteComment
*/
func (controller *Controller) DeleteCommentREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostCommentDelete{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentDeleteBind.With(err))
	}

	if err := controller.validateDeleteComment(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentDeleteValidation.With(err))
	}

	if err := controller._deleteComment(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}
