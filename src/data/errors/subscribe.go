package errors

import "tok-core/src/data/models"

var (
	SubscriptionExist = &models.Error{
		TechMessage:     "Subscription already exist",
		BusinessMessage: "Произошла ошибка",
	}
	SubscriptionNotExist = &models.Error{
		TechMessage:     "Already unsubscribed",
		BusinessMessage: "Произошла ошибка",
	}

	SubscribeOfProfileBind = &models.Error{
		TechMessage:     "Bind subscribe to user of profile model error",
		BusinessMessage: "Произошла ошибка",
	}

	SubscribeOfProfileValidate = &models.Error{
		TechMessage:     "Validate subscribe to user of profile error",
		BusinessMessage: "Произошла ошибка",
	}
	UnsubscribeOfProfileValidate = &models.Error{
		TechMessage:     "Validate unsubscribe to user of profile error",
		BusinessMessage: "Произошла ошибка",
	}

	SubscribeOfProfile = &models.Error{
		TechMessage:     "Subscribe to user of profile error",
		BusinessMessage: "Произошла ошибка",
	}
	UnsubscribeOfProfile = &models.Error{
		TechMessage:     "Unsubscribe to user of profile error",
		BusinessMessage: "Произошла ошибка",
	}
	SubscribersGet = &models.Error{
		TechMessage:     "Get profile subscribers error",
		BusinessMessage: "Произошла ошибка",
	}
	SubscribersExist = &models.Error{
		TechMessage:     "Getting subscription exist error",
		BusinessMessage: "Произошла ошибка",
	}
	SubscriptionsGet = &models.Error{
		TechMessage:     "Get profile subscriptions error",
		BusinessMessage: "Произошла ошибка",
	}
)
