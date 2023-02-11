package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Memory struct {
	Map
}

func (g Memory) incr(key string) (int64, error) {
	var value int64
	if err := g.Get(key, &value); err != nil {
		return value, fmt.Errorf("could not restore, %v = %v: %w", key, value, err)
	}
	if err := g.Set(key, value+1); err != nil {
		return value, fmt.Errorf("could not store, %v = %v + 1: %w", key, value, err)
	}
	return value, nil

}

type Map interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}) error
}

type Database struct {
	client *redis.Client
}

func NewDatabase() (Database, error) {
	db := Database{redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})}
	_, err := db.client.Ping(ctx).Result()
	return db, err

}
func (db Database) Set(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	log.Printf("%v.set %v, %v, %v", db, ctx, key, p)
	_, err = db.client.Set(ctx, key, p, 0).Result()
	return err
}

func (db Database) Get(key string, dest interface{}) error {
	get := db.client.Get(ctx, key)
	value, err := get.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, dest)
}
