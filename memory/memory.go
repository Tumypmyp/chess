package memory

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	
)



var ctx = context.Background()

type Memory interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}) error
	Incr(key string) (int64, error)
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
func (db Database) Incr(key string) (int64, error) {
	val, err := db.client.Incr(ctx, key).Result()
	return val, err
}
