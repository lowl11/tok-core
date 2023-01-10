package feed_event

import (
	"strings"
	"time"
	"tok-core/src/data/models"
)

const (
	explorePrefix      = "explore_"
	unauthorizedPrefix = "unauthorized_"
)

var (
	exploreFields = []string{
		"text", "keys",
	}
)

func (event *Event) notMyAccount(username string) map[string]any {
	return map[string]any{
		"term": map[string]any{
			"author": username,
		},
	}
}

func (event *Event) getExplore(date time.Time, username string, keys []string, page int) ([]models.PostElasticGet, error) {
	indexName := explorePrefix + date.Format("02-01-2006")

	var from int
	var size int

	if page == -1 {
		from = 0
		size = 1000
	} else {
		from = (page - 1) * 10
		size = 10
	}

	results, err := event.search.
		MultiMatch(indexName, strings.Join(keys, " "), exploreFields).
		Not(event.notMyAccount(username)).
		From(from).
		Size(size).
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

func (event *Event) getUnauthorized(date time.Time, username string, keys []string, page int) ([]models.PostElasticGet, error) {
	indexName := unauthorizedPrefix + date.Format("02-01-2006")

	var from int
	var size int

	if page == -1 {
		from = 0
		size = 1000
	} else {
		from = (page - 1) * 10
		size = 10
	}

	results, err := event.search.
		MultiMatch(indexName, strings.Join(keys, " "), exploreFields).
		Not(event.notMyAccount(username)).
		From(from).
		Size(size).
		Search()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		// ищем во вчерашнем индексе
		yesterdayIndexName := unauthorizedPrefix + time.Now().AddDate(0, 0, -1).Format("02-01-2006")

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
