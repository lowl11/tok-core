package notification_controller

import (
	"github.com/lowl11/lazy-collection/array"
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/data/entities"
	"tok-core/src/data/errors"
	"tok-core/src/data/models"
	"tok-core/src/definition"
)

func (controller *Controller) _read(username string, keys []string) *models.Error {
	logger := definition.Logger

	if keys == nil || len(keys) == 0 {
		return nil
	}

	if err := controller.notificationRepo.ReadItemList(username, keys); err != nil {
		logger.Error(err, "Read notification actions error", layers.Mongo)
		return errors.NotificationRead.With(err)
	}

	return nil
}

func (controller *Controller) _getInfo(username string) (*models.NotificationGetInfo, *models.Error) {
	logger := definition.Logger

	info, err := controller.notificationRepo.GetInfo(username)
	if err != nil {
		logger.Error(err, "Get notifications info error", layers.Mongo)
		return nil, errors.NotificationGetInfo.With(err)
	}

	usernames := make([]string, 0, len(info.Actions))
	array.NewWithList[entities.NotificationAction](info.Actions...).Each(func(item entities.NotificationAction) {
		usernames = append(usernames, item.Username)
	})

	userDynamic, err := controller.userRepo.GetDynamicByUsernames(usernames)
	if err != nil {
		return nil, errors.UserDynamicGet.With(err)
	}
	userList := type_list.NewWithList[entities.UserDynamicGet, models.UserDynamicGet](userDynamic...).Select(func(item entities.UserDynamicGet) models.UserDynamicGet {
		return models.UserDynamicGet{
			Username: item.Username,
			Name:     item.Name,
			Avatar:   item.Avatar,
		}
	})

	return &models.NotificationGetInfo{
		Username:        info.Username,
		NewActionsCount: info.NewActionsCount,
		Actions: type_list.NewWithList[entities.NotificationAction, models.NotificationAction](info.Actions...).Select(func(item entities.NotificationAction) models.NotificationAction {
			return models.NotificationAction{
				Status:     item.Status,
				User:       userList.Single(func(userInfo models.UserDynamicGet) bool { return userInfo.Username == item.Username }),
				ActionKey:  item.ActionKey,
				ActionCode: item.ActionCode,
				ActionBody: item.ActionBody,
				CreatedAt:  item.CreatedAt,
			}
		}).Slice(),
	}, nil
}

func (controller *Controller) _getCount(username string) (*models.NotificationGetCount, *models.Error) {
	logger := definition.Logger

	count, err := controller.notificationRepo.GetCount(username)
	if err != nil {
		logger.Error(err, "Get notifications count error", layers.Mongo)
		return nil, errors.NotificationGetCount.With(err)
	}

	if count == nil {
		logger.Warn("Notification count not found", layers.Mongo)
		return &models.NotificationGetCount{
			Username:        username,
			NewActionsCount: 0,
		}, nil
	}

	return &models.NotificationGetCount{
		Username:        count.Username,
		NewActionsCount: count.NewActionsCount,
	}, nil
}
