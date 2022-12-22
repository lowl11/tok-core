package controller

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"tok-core/src/data/models"
)

type Base struct {
	//
}

func (controller *Base) Error(ctx echo.Context, err *models.Error, status ...int) error {
	responseStatus := http.StatusInternalServerError
	if len(status) > 0 {
		responseStatus = status[0]
	}

	errorObject := &models.Response{
		Status:       "ERROR",
		Message:      err.BusinessMessage,
		InnerMessage: err.TechMessage,
	}
	return ctx.JSON(responseStatus, errorObject)
}

func (controller *Base) NotFound(ctx echo.Context, err *models.Error) error {
	errorObject := &models.Response{
		Status:       "ERROR",
		Message:      err.BusinessMessage,
		InnerMessage: err.TechMessage,
	}
	return ctx.JSON(http.StatusNotFound, errorObject)
}

func (controller *Base) Unauthorized(ctx echo.Context, err *models.Error) error {
	errorObject := &models.Response{
		Status:       "ERROR",
		Message:      err.BusinessMessage,
		InnerMessage: err.TechMessage,
	}
	return ctx.JSON(http.StatusUnauthorized, errorObject)
}

func (controller *Base) Ok(ctx echo.Context, response interface{}, messages ...string) error {
	defaultMessage := "OK"
	if len(messages) > 0 {
		defaultMessage = messages[0]
	}

	successObject := &models.Response{
		Status:       "OK",
		Message:      defaultMessage,
		InnerMessage: defaultMessage,
		Body:         response,
	}
	return ctx.JSON(http.StatusOK, successObject)
}

func (controller *Base) RequiredField(value interface{}, name string) error {
	if value == nil {
		return errors.New(fmt.Sprintf("Field %s is null, but it's required", name))
	}

	_, isString := value.(string)
	if isString && value.(string) == "" {
		return errors.New(fmt.Sprintf("Field %s is null or empty, but it's required", name))
	}

	_, isInt := value.(int)
	if isInt && value.(int) == 0 {
		return errors.New(fmt.Sprintf("Field %s is null or zero, but it's required", name))
	}

	return nil
}
