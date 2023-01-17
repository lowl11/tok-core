package post_repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"path/filepath"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (repo *Repository) Create(
	model *models.PostAddExtended,
	author, code string,
	uploadedPicturePath *string,
	customCategoryCode *string,
) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// если загрузили изображение
	var picturePath *string
	if model.Base.Picture != nil && uploadedPicturePath != nil {
		newPostPicture := "/images/post/" + author + "/post_" + code + "/post_picture" + filepath.Ext(*uploadedPicturePath)
		picturePath = &newPostPicture
	}

	var categoryCode string
	if customCategoryCode != nil {
		categoryCode = *customCategoryCode
	} else {
		categoryCode = model.Base.CategoryCode
	}

	// сущность БД
	entity := &entities.PostCreate{
		CategoryCode:   categoryCode,
		AuthorUsername: author,

		Text:    model.Base.Text,
		Picture: picturePath,
		Code:    code,

		PictureWidth:  model.ImageConfig.Width,
		PictureHeight: model.ImageConfig.Height,
	}

	// скрипт
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

func (repo *Repository) GetByUsername(username string, offset, size int) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// скрипт
	query := repo.Script("post", "get_by_username")

	rows, err := repo.connection.QueryxContext(ctx, query, username, offset, size)
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

func (repo *Repository) GetByUsernameList(usernameList []string, offset, size int) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// скрипт
	query := repo.Script("post", "get_by_username_list")

	rows, err := repo.connection.QueryxContext(ctx, query, pq.Array(usernameList), offset, size)
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

func (repo *Repository) GetByCategory(categoryCode string, offset, size int) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// скрипт
	query := repo.Script("post", "get_by_category")

	rows, err := repo.connection.QueryxContext(ctx, query, categoryCode, offset, size)
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

	// скрипт
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

func (repo *Repository) GetByCodeList(postCodes []string) ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// скрипт
	query := repo.Script("post", "get_by_code_list")

	rows, err := repo.connection.QueryxContext(ctx, query, pq.Array(postCodes))
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

func (repo *Repository) DeleteByCode(code string) error {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "delete_by_code")

	if err := repo.Transaction(repo.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, query, code); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetExplore() ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_explore")

	rows, err := repo.connection.QueryxContext(ctx, query)
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

func (repo *Repository) GetUnauthorized() ([]entities.PostGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "get_unauthorized")

	rows, err := repo.connection.QueryxContext(ctx, query)
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

func (repo *Repository) CountByUser(username string) (int, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post", "count_by_user")

	rows, err := repo.connection.QueryxContext(ctx, query, username)
	if err != nil {
		return 0, err
	}
	defer repo.CloseRows(rows)

	if rows.Next() {
		var count int
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}

	return 0, nil
}
