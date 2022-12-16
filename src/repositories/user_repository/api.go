package user_repository

import "tok-core/src/data/entities"

func (repository *Repository) GetByUsername(username string) (*entities.UserGet, error) {
	ctx, cancel := repository.Ctx()
	defer cancel()

	query := repository.Script("user", "get_by_username")

	rows, err := repository.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer repository.CloseRows(rows)

	if rows.Next() {
		user := entities.UserGet{}
		if err = rows.StructScan(&user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	return nil, nil
}
