package search_controller

import (
	"tok-core/src/data/models"
)

const (
	minQueryLength = 2
)

func (controller *Controller) validateUser(model *models.SearchUser) error {
	if err := controller.RequiredField(model.Query, "query"); err != nil {
		return err
	}

	return nil
}
