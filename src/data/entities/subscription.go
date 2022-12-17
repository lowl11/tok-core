package entities

type ProfileSubscribe struct {
	ProfileLogin   string `db:"profile_login"`
	SubscribeLogin string `db:"subscriber_login"`
}

type ProfileUnsubscribe struct {
	ProfileLogin   string `db:"profile_login"`
	SubscribeLogin string `db:"subscriber_login"`
}

type ProfileSubscriber struct {
	Username string `db:"username"`
}

type ProfileSubscription struct {
	Username string `db:"username"`
}
