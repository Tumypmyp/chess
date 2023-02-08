package main

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Memory interface {
	Get(key string, dest interface{}) error
	Set(key string, dest interface{}) error
}

type DataBase struct {
	client *redis.Client
}

func NewDataBase() DataBase {
	return DataBase{redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})}
}
func (db DataBase) Set(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = db.client.Set(ctx, key, p, 0).Result()
	return err
}

func (db DataBase) Get(key string, dest interface{}) error {
	get := db.client.Get(ctx, key)
	value, err := get.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(value, dest)
}
