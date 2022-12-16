package auth_controller

import (
	"errors"
	"tok-core/src/data/models"
)

func (controller *Controller) validateSignUp(model *models.Signup) error {
	if err := controller.RequiredField(model.Username, "username"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.Password, "password"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.RePassword, "re_password"); err != nil {
		return err
	}

	if model.Password != model.RePassword {
		return errors.New("passwords are not equal")
	}

	return nil
}

func (controller *Controller) validateLoginByCredentials(model *models.LoginByCredentials) error {
	if err := controller.RequiredField(model.Username, "username"); err != nil {
		return err
	}

	if err := controller.RequiredField(model.Password, "password"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateLoginByToken(model *models.LoginByToken) error {
	if err := controller.RequiredField(model.Token, "token"); err != nil {
		return err
	}

	return nil
}
