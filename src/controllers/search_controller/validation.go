package search_controller

import (
	"tok-core/src/data/models"
)

func (controller *Controller) validateUser(model *models.SearchUser) error {
	if err := controller.RequiredField(model.Query, "query"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateCategory(model *models.SearchCategory) error {
	if err := controller.RequiredField(model.Query, "query"); err != nil {
		return err
	}

	return nil
}

func (controller *Controller) validateSmart(model *models.SearchSmart) error {
	if err := controller.RequiredField(model.Query, "query"); err != nil {
		return err
	}

	return nil
}
