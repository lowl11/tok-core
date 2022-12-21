package user_repository

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (repository *Repository) Search(searchQuery string) ([]entities.UserGet, error) {
	ctx, cancel := repository.Ctx()
	defer cancel()

	query := repository.Script("user", "search")

	rows, err := repository.connection.QueryxContext(ctx, query, searchQuery)
	if err != nil {
		return nil, err
	}
	defer repository.CloseRows(rows)

	list := make([]entities.UserGet, 0)
	for rows.Next() {
		item := entities.UserGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

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

func (repository *Repository) UpdateProfile(username string, model *models.ProfileUpdate) error {
	ctx, cancel := repository.Ctx()
	defer cancel()

	// entity
	entity := &entities.ProfileUpdate{
		Username: username,
		Name:     model.Name,
		BIO:      model.BIO,
	}

	// query
	query := repository.Script("user", "update")

	if err := repository.Transaction(repository.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) UpdateAvatar(username, fileName string) error {
	ctx, cancel := repository.Ctx()
	defer cancel()

	// entity
	entity := &entities.ProfileAvatarUpdate{
		Username: username,
		Path:     "/images/profile/" + username + "/" + fileName,
	}

	// query
	query := repository.Script("user", "update_avatar")

	if err := repository.Transaction(repository.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) UpdateWallpaper(username, fileName string) error {
	ctx, cancel := repository.Ctx()
	defer cancel()

	// entity
	entity := &entities.ProfileWallpaperUpdate{
		Username: username,
		Path:     "/images/profile/" + username + "/" + fileName,
	}

	// query
	query := repository.Script("user", "update_wallpaper")

	if err := repository.Transaction(repository.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repository *Repository) UpdateContact(model *models.ProfileUpdateContact) error {
	return nil
}
