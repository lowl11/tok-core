package feed_event

import (
	"strings"
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

func (event *Event) GetExplore(username string, keys []string, page int) ([]models.PostElasticGet, error) {
	indexName := explorePrefix + time.Now().Format("02-01-2006")

	results, err := event.search.
		MultiMatch(indexName, strings.Join(keys, " "), exploreFields).
		Not(event.notMyAccount(username)).
		Search()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		// ищем во вчерашнем индексе
		yesterdayIndexName := explorePrefix + time.Now().AddDate(0, 0, -1).Format("02-01-2006")

		yesterdayResults, err := event.search.
			MultiMatch(yesterdayIndexName, strings.Join(keys, " "), exploreFields).
			Not(event.notMyAccount(username)).
			Search()
		if err != nil {
			return nil, err
		}

		return yesterdayResults, nil
	}

	return results, nil
}
