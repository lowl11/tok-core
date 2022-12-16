package client_session_event

import (
	"context"
	"time"
)

const (
	sessionPrefix = "cs_"
)

func (event *Event) ctx() (context.Context, func()) {
	return context.WithTimeout(context.Background(), time.Second*2)
}
