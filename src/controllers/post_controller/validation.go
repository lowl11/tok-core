package post_controller

import "tok-core/src/data/models"

func (controller *Controller) validateAddCategory(model *models.PostCategoryAdd) error {
	if err := controller.RequiredField(model.Name, "name"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validatePostCreate(model *models.PostAdd) error {
	if model.CustomCategory == nil {
		if err := controller.RequiredField(model.CategoryCode, "category"); err != nil {
			return err
		}
	} else {
		if err := controller.RequiredField(model.CustomCategory, "custom_category"); err != nil {
			return err
		}
	}

	if err := controller.RequiredField(model.Text, "text"); err != nil {
		return err
	}

	if model.Picture != nil {
		if err := controller.RequiredField(model.Picture.Name, "picture.name"); err != nil {
			return err
		}

		if err := controller.RequiredField(model.Picture.Buffer, "picture.buffer"); err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) validateLike(model *models.PostLike) error {
	if err := controller.RequiredField(model.PostCode, "post_code"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.PostCategory, "post_category"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateUnlike(model *models.PostUnlike) error {
	if err := controller.RequiredField(model.PostCode, "post_code"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.PostCategory, "post_category"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateAddComment(model *models.PostCommentAdd) error {
	if err := controller.RequiredField(model.PostAuthor, "post_author"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.PostCode, "post_code"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.CommentText, "comment_text"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateDeleteComment(model *models.PostCommentDelete) error {
	if err := controller.RequiredField(model.CommentCode, "comment_code"); err != nil {
		return err
	}

	if model.SubComment {
		if err := controller.RequiredField(model.ParentCommentCode, "parent_comment_code"); err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) validateLikeComment(model *models.PostCommentLike) error {
	if err := controller.RequiredField(model.CommentCode, "comment_code"); err != nil {
		return err
	}

	if model.SubComment {
		if err := controller.RequiredField(model.ParentCommentCode, "parent_comment_code"); err != nil {
			return err
		}
	}

	return nil
}

func (controller *Controller) validateUnlikeComment(model *models.PostCommentUnlike) error {
	if err := controller.RequiredField(model.CommentCode, "comment_code"); err != nil {
		return err
	}

	if model.SubComment {
		if err := controller.RequiredField(model.ParentCommentCode, "parent_comment_code"); err != nil {
			return err
		}
	}

	return nil
}
