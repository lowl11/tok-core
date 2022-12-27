package feed_event

import (
	"time"
	"tok-core/src/data/models"
)

func (event *Event) AddExplore(post *models.PostElasticAdd) error {
	// explore_27-12-2022
	indexName := explorePrefix + time.Now().Format("02-01-2006")

	// создать индекс перед записью (с проверкой на существование)
	if err := event.client.CreateIndex(indexName, nil); err != nil {
		return err
	}

	// создать запись
	return event.client.Insert(post.Code, indexName, post)
}

func (event *Event) DeleteExplore(postCode string) error {
	// explore_27-12-2022
	indexName := explorePrefix + time.Now().Format("02-01-2006")

	// удалить запись с индекса
	return event.client.DeleteItem(indexName, postCode)
}
