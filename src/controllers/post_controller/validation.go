package post_controller

import "tok-core/src/data/models"

func (controller *Controller) validatePostCreate(model *models.PostAdd) error {
	if err := controller.RequiredField(model.CategoryCode, "category"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.Text, "text"); err != nil {
		return err
	}

	return nil
}
