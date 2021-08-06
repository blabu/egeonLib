package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Model struct {
	expireTime time.Duration
	storage    *redis.Client
}

func GetCachedDB(ip, login, pass string, dbNumb int, expire int) (Model, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     ip,
		Username: login,
		Password: pass,
		DB:       dbNumb,
	})
	m := Model{
		expireTime: time.Duration(expire) * time.Second,
		storage:    client,
	}
	err := client.Ping(context.Background()).Err()
	if err != nil {
		m.storage = nil
	}
	return m, err
}

func (m Model) Close() error {
	if m.storage != nil {
		m.storage.Save(context.Background())
		return m.storage.Close()
	}
	return nil
}

func (m Model) Save(ctx context.Context) {
	if m.storage != nil {
		m.storage.BgSave(ctx)
	}
}

func (m Model) Set(ctx context.Context, key string, resp *Responce) error {
	if resp == nil {
		return errors.New("bad responce for cache")
	}
	if data, err := json.Marshal(*resp); err == nil {
		m.storage.Set(ctx, key, data, m.expireTime)
	} else {
		return err
	}
	return nil
}

func (m Model) Get(ctx context.Context, key string) (*Responce, error) {
	data, err := m.storage.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var resp Responce
	err = json.Unmarshal(data, &resp)
	return &resp, err
}
