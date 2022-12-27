package feed_event

import (
	"github.com/lowl11/lazy-elastic/elastic_client"
	"github.com/lowl11/lazy-elastic/elastic_search"
	"github.com/lowl11/lazy-elastic/es_api"
	"tok-core/src/data/models"
)

type Event struct {
	client *elastic_client.Event
	search *elastic_search.Event[models.PostGet]
}

func Create(servers []string, _, _ string) *Event {
	baseURL := servers[0]

	return &Event{
		client: es_api.NewClient(baseURL),
		search: es_api.NewSearch[models.PostGet](baseURL),
	}
}
