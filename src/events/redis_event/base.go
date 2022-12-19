package redis_event

import (
	"context"
	"time"
)

type Base struct {
	//
}

func (event *Base) Ctx() (context.Context, func()) {
	return context.WithTimeout(context.Background(), time.Second*2)
}
