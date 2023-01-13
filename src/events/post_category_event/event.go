package post_category_event

import (
	"github.com/lowl11/lazy-collection/array"
	"sync"
	"tok-core/src/data/models"
)

type Event struct {
	categories *array.Array[models.PostCategoryGet]

	mutex sync.Mutex
}

func Create() *Event {
	return &Event{
		categories: array.New[models.PostCategoryGet](),
		mutex:      sync.Mutex{},
	}
}
