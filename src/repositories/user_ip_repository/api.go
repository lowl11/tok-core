package user_ip_repository

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/data/entities"
)

func (repo *Repository) Get(username string) ([]string, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("user_ip", "get")

	rows, err := repo.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]string, 0)
	for rows.Next() {
		var ipAddress string
		if err = rows.Scan(&ipAddress); err != nil {
			return nil, err
		}
		list = append(list, ipAddress)
	}

	return list, nil
}

func (repo *Repository) New(username, ipAddress string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// сущность
	entity := &entities.UserIpNew{
		Username:  username,
		IpAddress: ipAddress,
	}

	// запрос
	query := repo.Script("user_ip", "new")

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

func (repo *Repository) Has(username, ipAddress string) (bool, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("user_ip", "has")

	rows, err := repo.connection.QueryxContext(ctx, query, username, ipAddress)
	if err != nil {
		return false, err
	}
	defer repo.CloseRows(rows)

	return rows.Next(), nil
}
