package notification_controller

import (
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) _getInfo(username string) (*models.NotificationGetInfo, *models.Error) {
	logger := definition.Logger

	info, err := controller.notificationRepo.GetInfo(username)
	if err != nil {
		logger.Error(err, "Get notifications info error", layers.Mongo)
		return nil, errors.NotificationGetInfo.With(err)
	}

	return &models.NotificationGetInfo{
		Username:        info.Username,
		NewActionsCount: info.NewActionsCount,
	}, nil
}

func (controller *Controller) _getCount(username string) (*models.NotificationGetCount, *models.Error) {
	logger := definition.Logger

	count, err := controller.notificationRepo.GetCount(username)
	if err != nil {
		logger.Error(err, "Get notifications count error", layers.Mongo)
		return nil, errors.NotificationGetCount.With(err)
	}

	return &models.NotificationGetCount{
		Username:        count.Username,
		NewActionsCount: count.NewActionsCount,
	}, nil
}
