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

func (controller *Controller) validateUploadAvatar(model *models.ImageAvatar) error {
	if err := controller.RequiredField(model.Name, "name"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.Buffer, "buffer"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateUploadWallpaper(model *models.ImageWallpaper) error {
	if err := controller.RequiredField(model.Name, "name"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.Buffer, "buffer"); err != nil {
		return err
	}

	return nil
}
