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

func (controller *Controller) _read(username string) *models.Error {
	logger := definition.Logger

	// get keys
	notifications, err := controller.notificationRepo.GetUnreadInfo(username)
	if err != nil {
		logger.Error(err, "Get unread notifications error", layers.Mongo)
		return errors.NotificationGetInfo.With(err)
	}

	keys := type_list.NewWithList[entities.NotificationGet, string](notifications...).Select(func(item entities.NotificationGet) string {
		return item.ActionKey
	}).Slice()

	if err = controller.notificationRepo.ReadItemList(username, keys); err != nil {
		logger.Error(err, "Read notification actions error", layers.Mongo)
		return errors.NotificationRead.With(err)
	}

	return nil
}

func (controller *Controller) _getInfo(username string, page int) ([]models.NotificationGet, *models.Error) {
	logger := definition.Logger

	from := (page - 1) * 10

	notifications, err := controller.notificationRepo.GetInfo(username, from)
	if err != nil {
		logger.Error(err, "Get notifications info error", layers.Mongo)
		return nil, errors.NotificationGetInfo.With(err)
	}

	notificationsLength := len(notifications)
	usernames := make([]string, 0, notificationsLength)
	postCodes := make([]string, 0, notificationsLength)

	array.NewWithList[entities.NotificationGet](notifications...).Each(func(item entities.NotificationGet) {
		usernames = append(usernames, item.ActionAuthor)
		if item.ActionBody != nil && item.ActionBody.PostCode != nil {
			postCodes = append(postCodes, *item.ActionBody.PostCode)
		}
	})

	userDynamic, err := controller.userRepo.GetDynamicByUsernames(usernames)
	if err != nil {
		logger.Error(err, "Get users by usernames error", layers.Mongo)
		return nil, errors.UserDynamicGet.With(err)
	}
	userList := type_list.NewWithList[entities.UserDynamicGet, models.UserDynamicGet](userDynamic...).Select(func(item entities.UserDynamicGet) models.UserDynamicGet {
		return models.UserDynamicGet{
			Username: item.Username,
			Name:     item.Name,
			Avatar:   item.Avatar,
		}
	})

	posts, err := controller.postRepo.GetByCodeList(postCodes)
	if err != nil {
		logger.Error(err, "Get posts by codes error", layers.Mongo)
		return nil, errors.PostsGetByCode.With(err)
	}

	postList := type_list.NewWithList[entities.PostGet, models.NotificationPostGet](posts...).Select(func(item entities.PostGet) models.NotificationPostGet {
		return models.NotificationPostGet{
			Code:  item.Code,
			Image: item.Picture,
			Text:  item.Text,
		}
	})

	go func() {
		if err = controller._read(username); err != nil {
			return
		}
	}()

	return type_list.NewWithList[entities.NotificationGet, models.NotificationGet](notifications...).Select(func(item entities.NotificationGet) models.NotificationGet {
		var post *models.NotificationPostGet
		if item.ActionBody != nil && item.ActionBody.PostCode != nil {
			post = postList.Single(func(postItem models.NotificationPostGet) bool { return postItem.Code == *item.ActionBody.PostCode })
		}

		return models.NotificationGet{
			Status:       item.Status,
			User:         userList.Single(func(userItem models.UserDynamicGet) bool { return userItem.Username == item.ActionAuthor }),
			ActionAuthor: item.ActionAuthor,
			ActionKey:    item.ActionKey,
			ActionCode:   item.ActionCode,
			ActionBody: &models.NotificationBody{
				Post: post,
			},
			CreatedAt: item.CreatedAt,
		}
	}).Slice(), nil
}

func (controller *Controller) _getCount(username string) (int, *models.Error) {
	logger := definition.Logger

	count, err := controller.notificationRepo.GetCount(username)
	if err != nil {
		logger.Error(err, "Get notifications count error", layers.Mongo)
		return 0, errors.NotificationGetCount.With(err)
	}

	return count, nil
}
