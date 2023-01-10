package feed_event

import (
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazy-elastic/es_model"
	"time"
	"tok-core/src/data/models"
)

const (
	AnonymousUsername = "L1RtsIGA1J5MGWcc"
)

func (event *Event) ExploreTodayExist() bool {
	return event.client.ExistIndex(explorePrefix + time.Now().Format("02-01-2006"))
}

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

func (event *Event) AddListExplore(date time.Time, posts []models.PostElasticAdd) error {
	// explore_27-12-2022
	indexName := explorePrefix + date.Format("02-01-2006")

	// создать индекс перед записью (с проверкой на существование)
	if err := event.client.CreateIndex(indexName, nil); err != nil {
		return err
	}

	// создать запись
	return event.client.InsertMultiple(indexName, type_list.NewWithList[models.PostElasticAdd, es_model.InsertMultipleData](posts...).Select(func(item models.PostElasticAdd) es_model.InsertMultipleData {
		return es_model.InsertMultipleData{
			ID:     item.Code,
			Object: item,
		}
	}).Slice())
}

func (event *Event) AddListUnauthorized(date time.Time, posts []models.PostElasticAdd) error {
	// explore_27-12-2022
	indexName := explorePrefix + date.Format("02-01-2006")

	// создать индекс перед записью (с проверкой на существование)
	if err := event.client.CreateIndex(indexName, nil); err != nil {
		return err
	}

	// создать запись
	return event.client.InsertMultiple(indexName, type_list.NewWithList[models.PostElasticAdd, es_model.InsertMultipleData](posts...).Select(func(item models.PostElasticAdd) es_model.InsertMultipleData {
		return es_model.InsertMultipleData{
			ID:     item.Code,
			Object: item,
		}
	}).Slice())
}

func (event *Event) DeletePostExplore(postCode string) error {
	// explore_27-12-2022
	indexName := explorePrefix + time.Now().Format("02-01-2006")

	// удалить запись с индекса
	return event.client.DeleteItem(indexName, postCode)
}

func (event *Event) GetExploreToday(username string, keys []string, page int) ([]models.PostElasticGet, error) {
	return event.getExplore(time.Now(), username, keys, page)
}

func (event *Event) GetExploreTomorrow(username string, keys []string, page int) ([]models.PostElasticGet, error) {
	return event.getExplore(time.Now().AddDate(0, 0, 1), username, keys, page)
}

func (event *Event) FillExploreToday(posts []models.PostElasticAdd) error {
	return event.AddListExplore(time.Now(), posts)
}

func (event *Event) FillExploreTomorrow(posts []models.PostElasticAdd) error {
	return event.AddListExplore(time.Now().AddDate(0, 0, 1), posts)
}

func (event *Event) UnauthorizedTodayExist() bool {
	return event.client.ExistIndex(unauthorizedPrefix + time.Now().Format("02-01-2006"))
}

func (event *Event) GetUnauthorizedToday(username string, keys []string, page int) ([]models.PostElasticGet, error) {
	return event.getUnauthorized(time.Now(), username, keys, page)
}

func (event *Event) GetUnauthorizedTomorrow(username string, keys []string, page int) ([]models.PostElasticGet, error) {
	return event.getUnauthorized(time.Now().AddDate(0, 0, 1), username, keys, page)
}

func (event *Event) FillUnauthorizedToday(posts []models.PostElasticAdd) error {
	return event.AddListUnauthorized(time.Now(), posts)
}

func (event *Event) FillUnauthorizedTomorrow(posts []models.PostElasticAdd) error {
	return event.AddListUnauthorized(time.Now().AddDate(0, 0, 1), posts)
}

func (event *Event) DeletePostUnauthorized(postCode string) error {
	// explore_27-12-2022
	indexName := unauthorizedPrefix + time.Now().Format("02-01-2006")

	// удалить запись с индекса
	return event.client.DeleteItem(indexName, postCode)
}
