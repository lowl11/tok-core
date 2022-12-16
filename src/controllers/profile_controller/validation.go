package profile_controller

import "tok-core/src/data/models"

func (controller *Controller) validateUpdate(model *models.ProfileUpdate) error {
	if err := controller.RequiredField(model.Name, "name"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.BIO, "bio"); err != nil {
		return err
	}

	return nil
}
