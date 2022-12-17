package client_session_event

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

func (event *Event) Delete(token, username string) error {
	ctx, cancel := event.ctx()
	defer cancel()

	return event.client.Del(ctx, sessionPrefix+username+"_"+token).Err()
}

func (event *Event) DeleteByToken(token string) error {
	ctx, cancel := event.ctx()
	defer cancel()

	keys, err := event.client.Keys(ctx, sessionPrefix+"*_"+token).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return event.client.Del(ctx, keys[0]).Err()
}

func (event *Event) DeleteByUsername(username string) error {
	ctx, cancel := event.ctx()
	defer cancel()

	keys, err := event.client.Keys(ctx, sessionPrefix+username+"_*").Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return event.client.Del(ctx, keys[0]).Err()
}

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
	if err = event.client.Set(ctx, sessionPrefix+session.Username+"_"+token, entityInBytes, time.Hour*24*7).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func (event *Event) Get(token, username string) (*entities.ClientSession, error) {
	ctx, cancel := event.ctx()
	defer cancel()

	// получить сессию в байтах
	sessionInBytes, err := event.client.Get(ctx, sessionPrefix+username+"_"+token).Bytes()
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

func (event *Event) GetByUsername(username string) (*entities.ClientSession, error) {
	ctx, cancel := event.ctx()
	defer cancel()

	keys, err := event.client.Keys(ctx, sessionPrefix+username+"_*").Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, errors.New("session not found")
	}

	// получить сессию в байтах
	sessionInBytes, err := event.client.Get(ctx, keys[0]).Bytes()
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

func (event *Event) GetByToken(token string) (*entities.ClientSession, error) {
	ctx, cancel := event.ctx()
	defer cancel()

	// найти ключ
	keys, err := event.client.Keys(ctx, sessionPrefix+"*_"+token).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, errors.New("session not found")
	}

	// получить сессию в байтах
	sessionInBytes, err := event.client.Get(ctx, keys[0]).Bytes()
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

func (event *Event) Update(session *entities.ClientSession, token string) error {
	ctx, cancel := event.ctx()
	defer cancel()

	sessionInBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	if err = event.client.Set(ctx, sessionPrefix+session.Username+"_"+token, sessionInBytes, time.Hour*24*7).Err(); err != nil {
		return err
	}

	return nil
}
