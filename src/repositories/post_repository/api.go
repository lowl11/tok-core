package post_repository

import (
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (repo *Repository) Create(model *models.PostAdd, author string) (int, error) {
	return 0, nil
}

func (repo *Repository) GetByUser(username string) ([]entities.PostGet, error) {
	return nil, nil
}
