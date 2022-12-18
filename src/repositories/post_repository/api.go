package post_repository

import (
	"github.com/jmoiron/sqlx"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (repo *Repository) Create(model *models.PostAdd, author, code string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// entity
	entity := &entities.PostCreate{
		CategoryCode:   model.CategoryCode,
		AuthorUsername: author,

		Text:    model.Text,
		Picture: model.Picture,
		Code:    code,
	}

	// query
	query := repo.Script("post", "create")

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

func (repo *Repository) GetByUsername(username string) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_by_username")

	rows, err := repo.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.PostGet, 0)
	for rows.Next() {
		item := entities.PostGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) GetByUsernameList(usernameList []string) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_by_username_list")

	rows, err := repo.connection.QueryxContext(ctx, query, usernameList)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.PostGet, 0)
	for rows.Next() {
		item := entities.PostGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) GetByCategory(categoryCode string) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_by_category")

	rows, err := repo.connection.QueryxContext(ctx, query, categoryCode)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.PostGet, 0)
	for rows.Next() {
		item := entities.PostGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) GetByCode(code string) (*entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_by_code")

	rows, err := repo.connection.QueryxContext(ctx, query, code)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	if rows.Next() {
		item := entities.PostGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		return &item, nil
	}

	return nil, nil
}