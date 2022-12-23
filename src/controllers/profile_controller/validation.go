package profile_controller

import "tok-core/src/data/models"

func (controller *Controller) validateSubscribeProfile(model *models.ProfileSubscribe) error {
	if err := controller.RequiredField(model.Username, "username"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateUnsubscribeProfile(model *models.ProfileUnsubscribe) error {
	if err := controller.RequiredField(model.Username, "username"); err != nil {
		return err
	}

	return nil
}

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

func (controller *Controller) validateUpdateContacts(_ *models.ProfileUpdateContact) error {
	return nil
}
