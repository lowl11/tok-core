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
	Username string  `db:"username"`
	Name     *string `db:"name"`
	Avatar   *string `db:"avatar"`
}

type ProfileSubscription struct {
	Username string  `db:"username"`
	Name     *string `db:"name"`
	Avatar   *string `db:"avatar"`
}
