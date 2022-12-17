package subscription_repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lowl11/lazy-collection/type_list"
	"tok-core/src/data/entities"
)

func (repo *Repository) ProfileSubscribe(profileLogin, subscriberLogin string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// entity
	entity := &entities.ProfileSubscribe{
		ProfileLogin:   profileLogin,
		SubscribeLogin: subscriberLogin,
	}

	// query
	query := repo.Script("subscribe", "profile")

	if err := repo.Transaction(repo.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ProfileSubscribers(username string) (*type_list.TypeList[entities.ProfileSubscriber, string], error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("subscribe", "profile_subscribers")

	rows, err := repo.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.ProfileSubscriber, 0)
	for rows.Next() {
		item := entities.ProfileSubscriber{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return type_list.NewWithList[entities.ProfileSubscriber, string](list...), nil
}

func (repo *Repository) ProfileSubscribersCount(username string) (int, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("subscribe", "profile_subscribers_count")

	var count int
	if err := repo.connection.QueryRowxContext(ctx, query, username).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *Repository) ProfileSubscriptions(username string) (*type_list.TypeList[entities.ProfileSubscription, string], error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("subscribe", "profile_subscriptions")

	rows, err := repo.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.ProfileSubscription, 0)
	for rows.Next() {
		item := entities.ProfileSubscription{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return type_list.NewWithList[entities.ProfileSubscription, string](list...), nil
}

func (repo *Repository) ProfileSubscriptionsCount(username string) (int, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("subscribe", "profile_subscriptions_count")

	var count int
	if err := repo.connection.QueryRowxContext(ctx, query, username).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
