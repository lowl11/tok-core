package entities

type ProfileSubscribe struct {
	ProfileUsername   string `db:"profile_username"`
	SubscribeUsername string `db:"subscriber_username"`
}

type ProfileUnsubscribe struct {
	ProfileUsername   string `db:"profile_username"`
	SubscribeUsername string `db:"subscriber_username"`
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
