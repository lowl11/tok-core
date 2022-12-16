package client_session_event

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (event *Event) Create(session *models.ClientSessionCreate) (string, error) {
	ctx, cancel := event.ctx()
	defer cancel()

	// создаем токен
	token := uuid.New().String()

	// сбор
	entity := &entities.ClientSession{
		Username:  session.Username,
		Name:      session.Name,
		BIO:       session.BIO,
		Avatar:    session.Avatar,
		Wallpaper: session.Wallpaper,
	}

	// превращаем сессию в байты
	entityInBytes, err := json.Marshal(entity)
	if err != nil {
		return "", err
	}

	// записываем в redis
	if err = event.client.Set(ctx, sessionPrefix+token, entityInBytes, time.Hour*24*7).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func (event *Event) Get(token string) (*entities.ClientSession, error) {
	ctx, cancel := event.ctx()
	defer cancel()

	// получить сессию в байтах
	sessionInBytes, err := event.client.Get(ctx, sessionPrefix+token).Bytes()
	if err != nil {
		return nil, err
	}

	// вдруг записался NULL
	if sessionInBytes == nil {
		return nil, errors.New("session not found")
	}

	// парсим байты в структуру
	entityGet := entities.ClientSession{}
	if err = json.Unmarshal(sessionInBytes, &entityGet); err != nil {
		return nil, err
	}

	return &entityGet, nil
}
