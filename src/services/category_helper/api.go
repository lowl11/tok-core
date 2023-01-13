package category_helper

import (
	"github.com/lowl11/lazy-collection/type_list"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
	"tok-core/src/events/post_category_event"
	"tok-core/src/repositories/post_category_repository"
)

func Load(postCategory *post_category_event.Event, categoryRepo *post_category_repository.Repository) {
	categories, err := categoryRepo.GetAll()
	if err != nil {
		//
	}

	postCategory.Load(
		type_list.NewWithList[entities.PostCategoryGet, models.PostCategoryGet](categories...).Select(func(item entities.PostCategoryGet) models.PostCategoryGet {
			return models.PostCategoryGet{
				Code: item.Code,
				Name: item.Name,
			}
		}).Slice(),
	)
}
