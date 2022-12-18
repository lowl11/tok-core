package post_category_repository

import "tok-core/src/data/entities"

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
