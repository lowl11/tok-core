package post_category_repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/mehanizm/iuliia-go"
	"strings"
	"tok-core/src/data/entities"
)

func (repo *Repository) GetAll() ([]entities.PostCategoryGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post_category", "get_all")

	rows, err := repo.connection.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.PostCategoryGet, 0)
	for rows.Next() {
		item := entities.PostCategoryGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) Search(searchQuery string) ([]entities.PostCategoryGet, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	query := repo.Script("post_category", "search")

	rows, err := repo.connection.QueryxContext(ctx, query, searchQuery)
	if err != nil {
		return nil, err
	}
	defer repo.CloseRows(rows)

	list := make([]entities.PostCategoryGet, 0)
	for rows.Next() {
		item := entities.PostCategoryGet{}
		if err = rows.StructScan(&item); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *Repository) Create(name string) (string, error) {
	ctx, cancel := repo.Ctx()
	defer cancel()

	// генерируем код категории на латинице
	// например Образование -> obrazovanie
	categoryCode := strings.ReplaceAll(strings.ToLower(iuliia.Wikipedia.Translate(name)), " ", "_")

	// запрос
	query := repo.Script("post_category", "create")

	// сущность
	entity := &entities.PostCategoryCreate{
		Name: name,
		Code: categoryCode,
	}

	if err := repo.Transaction(repo.connection, func(tx *sqlx.Tx) error {
		if _, err := tx.NamedExecContext(ctx, query, entity); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return categoryCode, nil
}
