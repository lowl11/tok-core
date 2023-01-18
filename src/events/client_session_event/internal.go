package client_session_event

import (
	"tok-core/src/data/entities"
)

const (
	sessionPrefix = "cs_"

	commonUsername = "u2lLaHKW"
	commonPassword = "fs21xmiqddddmcm1-vz2w3tvdqfa3r4no5qrc8_o6ku="
	commonToken    = "1a46181f-7f8b-4498-b230-27ac5ff7a288"
)

func (event *Event) getMockSession(token, username string) *entities.ClientSession {
	if token == commonToken && username == commonUsername {
		return &entities.ClientSession{
			ID:       777,
			Username: commonUsername,
			Password: commonPassword,
			Subscriptions: entities.ClientSessionSubscribes{
				SubscriptionCount: 0,
				SubscriberCount:   0,
				Subscriptions:     make([]string, 0),
				Subscribers:       make([]string, 0),
			},
		}
	}

	return nil
}
