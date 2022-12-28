package feed_event

import (
	"github.com/lowl11/lazy-collection/type_list"
	"github.com/lowl11/lazy-elastic/es_model"
	"strings"
	"time"
	"tok-core/src/data/entities"
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

func (event *Event) AddListExplore(posts []models.PostElasticAdd) error {
	// explore_27-12-2022
	indexName := explorePrefix + time.Now().Format("02-01-2006")

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

func (event *Event) FillExplore(posts []entities.PostGet) error {
	postsList := type_list.NewWithList[entities.PostGet, models.PostElasticAdd](posts...)
	insertList := postsList.Select(func(item entities.PostGet) models.PostElasticAdd {
		pictureModel := &models.PostGetPicture{
			Path:   item.Picture,
			Width:  item.PictureWidth,
			Height: item.PictureHeight,
		}

		return models.PostElasticAdd{
			Code:         item.Code,
			Text:         item.Text,
			Category:     item.CategoryCode,
			CategoryName: item.CategoryName,
			Picture:      pictureModel,
			Author:       item.AuthorUsername,
			CreatedAt:    item.CreatedAt,

			Keys: []string{item.CategoryCode, item.CategoryName},
		}
	})

	return event.AddListExplore(insertList.Slice())
}
