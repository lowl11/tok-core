package feed_event

import (
	"time"
	"tok-core/src/data/models"
)

func (event *Event) AddRecommendation(post *models.PostElasticAdd) error {
	// recmd_27-12-2022
	indexName := recmdPrefix + time.Now().Format("02-01-2006")

	// создать индекс перед записью (с проверкой на существование)
	if err := event.client.CreateIndex(indexName, nil); err != nil {
		return err
	}

	// создать запись
	return event.client.Insert(post.Code, indexName, post)
}

func (event *Event) DeleteRecommendation(postCode string) error {
	// recmd_27-12-2022
	indexName := recmdPrefix + time.Now().Format("02-01-2006")

	// удалить запись с индекса
	return event.client.DeleteItem(indexName, postCode)
}
