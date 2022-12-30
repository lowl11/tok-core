package repository_mongo

import (
	"context"
	"github.com/lowl11/lazylog/layers"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"tok-core/src/definition"
)

type Base struct {
	//
}

func CreateBase() Base {
	return Base{}
}

func (repo *Base) Ctx(customTimeout ...time.Duration) (context.Context, func()) {
	defaultTimeout := time.Second * 5
	if len(customTimeout) > 0 {
		defaultTimeout = customTimeout[0]
	}
	return context.WithTimeout(context.Background(), defaultTimeout)
}

func (repo *Base) CloseCursor(cursor *mongo.Cursor) {
	ctx, cancel := repo.Ctx(time.Second * 2)
	defer cancel()

	if err := cursor.Close(ctx); err != nil {
		definition.Logger.Error(err, "Close cursor error", layers.Mongo)
	}
}

func (repo *Base) LogError(err error) {
	if err != nil {
		definition.Logger.Error(err, "Cursor error", layers.Mongo)
	}
}
