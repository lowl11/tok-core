package mongo_service

import (
	"go.mongodb.org/mongo-driver/bson"
)

type FilterService struct {
	isAnd      bool
	conditions map[string]any
}

func (service *FilterService) And() *FilterService {
	service.isAnd = true
	return service
}

func (service *FilterService) Or() *FilterService {
	service.isAnd = false
	return service
}

func (service *FilterService) Eq(key string, value any) *FilterService {
	service.conditions[key] = value
	return service
}

func (service *FilterService) Get() bson.M {
	build := bson.M{}
	buildConditions := service.getConditions()

	if service.isAnd {
		build = bson.M{
			"$and": buildConditions,
		}
	} else {
		build = bson.M{
			"$or": buildConditions,
		}
	}

	return build
}

func (service *FilterService) getConditions() []interface{} {
	buildConditions := make([]interface{}, 0, len(service.conditions))
	for key, value := range service.conditions {
		buildConditions = append(buildConditions, bson.M{key: value})
	}
	return buildConditions
}
