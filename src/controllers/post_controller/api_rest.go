package post_controller

import (
	"github.com/labstack/echo/v4"
	"strconv"
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

	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page == 0 {
		page = 1
	}

	likes, err := controller._getLikes(session, postCode, page)
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

/*
	LikeCommentREST обертка для _likeComment
*/
func (controller *Controller) LikeCommentREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostCommentLike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentLikeBind.With(err))
	}

	if err := controller.validateLikeComment(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentLikeValidation.With(err))
	}

	if err := controller._likeComment(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	UnlikeCommentREST обертка для _unlikeComment
*/
func (controller *Controller) UnlikeCommentREST(ctx echo.Context) error {
	session := ctx.Get("client_session").(*entities.ClientSession)

	model := models.PostCommentUnlike{}
	if err := ctx.Bind(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentUnlikeBind.With(err))
	}

	if err := controller.validateUnlikeComment(&model); err != nil {
		return controller.Error(ctx, errors.PostCommentUnlikeValidation.With(err))
	}

	if err := controller._unlikeComment(session, &model); err != nil {
		return controller.Error(ctx, err)
	}

	return controller.Ok(ctx, "OK")
}

/*
	FillUnauthorizedJob обертка для _fillUnauthorizedFeed
*/
func (controller *Controller) FillUnauthorizedJob() error {
	if err := controller._fillUnauthorizedFeed(); err != nil {
		return err
	}

	return nil
}

/*
	FillExploreJob обертка для _fillExploreFeed
*/
func (controller *Controller) FillExploreJob() error {
	if err := controller._fillExploreFeed(); err != nil {
		return err
	}

	return nil
}
