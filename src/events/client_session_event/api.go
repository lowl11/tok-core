package client_session_event

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
	"tok-core/src/data/entities"
	"tok-core/src/data/models"
)

/*
	Delete удаление сессии по токену и юзернейму
*/
func (event *Event) Delete(token, username string) error {
	ctx, cancel := event.Ctx()
	defer cancel()

	// удаление сессии
	return event.client.Del(ctx, sessionPrefix+username+"_"+token).Err()
}

/*
	DeleteByToken удаление сессии по токену
*/
func (event *Event) DeleteByToken(token string) error {
	ctx, cancel := event.Ctx()
	defer cancel()

	// сначала получаем совпадающие ключи
	keys, err := event.client.Keys(ctx, sessionPrefix+"*_"+token).Result()
	if err != nil {
		return err
	}

	// сессии нет если не нашли ключи
	if len(keys) == 0 {
		return nil
	}

	// удаление сессии
	return event.client.Del(ctx, keys[0]).Err()
}

/*
	DeleteByUsername удаление сессии по юзернейму
*/
func (event *Event) DeleteByUsername(username string) error {
	ctx, cancel := event.Ctx()
	defer cancel()

	// сначала получаем совпадающие ключи
	keys, err := event.client.Keys(ctx, sessionPrefix+username+"_*").Result()
	if err != nil {
		return err
	}

	// сессии нет если не нашли ключи
	if len(keys) == 0 {
		return nil
	}

	// удаление сессии
	return event.client.Del(ctx, keys[0]).Err()
}

/*
	Create создание сессии
*/
func (event *Event) Create(session *models.ClientSessionCreate) (string, error) {
	ctx, cancel := event.Ctx()
	defer cancel()

	// создаем токен
	token := uuid.New().String()

	// сбор сущности сессии
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

/*
	Get получение сессии по токену и юзернейму
*/
func (event *Event) Get(token, username string) (*entities.ClientSession, error) {
	ctx, cancel := event.Ctx()
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

/*
	GetByUsername получение сессии по юзернейму
*/
func (event *Event) GetByUsername(username string) (*entities.ClientSession, string, error) {
	ctx, cancel := event.Ctx()
	defer cancel()

	// находим подходящие ключи
	keys, err := event.client.Keys(ctx, sessionPrefix+username+"_*").Result()
	if err != nil {
		return nil, "", err
	}

	// ошибка если ключи не найдены (значит сессии нет)
	if len(keys) == 0 {
		return nil, "", errors.New("session not found")
	}

	// получить сессию в байтах
	sessionInBytes, err := event.client.Get(ctx, keys[0]).Bytes()
	if err != nil {
		return nil, "", err
	}

	// вдруг записался NULL
	if sessionInBytes == nil {
		return nil, "", errors.New("session not found")
	}

	// парсим байты в структуру
	entityGet := entities.ClientSession{}
	if err = json.Unmarshal(sessionInBytes, &entityGet); err != nil {
		return nil, "", err
	}

	// парсим токен
	token := keys[0]
	token = strings.ReplaceAll(token, sessionPrefix, "")
	token = strings.ReplaceAll(token, username, "")
	token = strings.ReplaceAll(token, "_", "")

	return &entityGet, token, nil
}

/*
	GetByToken получение сессии по токену
*/
func (event *Event) GetByToken(token string) (*entities.ClientSession, error) {
	ctx, cancel := event.Ctx()
	defer cancel()

	// ищем подходящие ключи
	keys, err := event.client.Keys(ctx, sessionPrefix+"*_"+token).Result()
	if err != nil {
		return nil, err
	}

	// ошибка если не нашли ключи (значит сессии нет)
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

/*
	Update обновление существующей сессии
*/
func (event *Event) Update(session *entities.ClientSession, token string) error {
	ctx, cancel := event.Ctx()
	defer cancel()

	// парсим сессию в байты
	sessionInBytes, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// обновляем сессию в redis
	if err = event.client.Set(ctx, sessionPrefix+session.Username+"_"+token, sessionInBytes, time.Hour*24*7).Err(); err != nil {
		return err
	}

	return nil
}
