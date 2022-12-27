package feed_event

import (
	"time"
	"tok-core/src/data/models"
)

func (event *Event) AddRecommendation(post *models.PostElasticAdd) error {
	// recmd_27-12-2022
	indexName := recmdPrefix + time.Now().Format("02-01-2006")

	// create index before recording (with exist check)
	if err := event.client.CreateIndex(indexName, nil); err != nil {
		return err
	}

	// create record
	return event.client.Insert(post.Code, indexName, post)
}

func (event *Event) CreateUserFeed() error {
	return nil
}

func (event *Event) CreateCategoryFeed() error {
	return nil
}

func (event *Event) UpdateUserFeed() error {
	return nil
}

func (event *Event) UpdateCategoryFeed() error {
	return nil
}
