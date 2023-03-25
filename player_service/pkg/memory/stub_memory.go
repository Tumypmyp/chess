package memory 

import (
	"encoding/json"
	"errors"
	"log"
)

func deepCopy(a, b interface{}) {
	byt, err := json.Marshal(a)
	if err != nil {
		log.Fatalf("%v", err)

	}
	json.Unmarshal(byt, b)
}

type StubDatabase struct {
	DB *map[string]interface{}
}

func NewStubDatabase() StubDatabase {
	db := make(map[string]interface{})
	return StubDatabase{DB: &db}
}

func (s StubDatabase) Get(key string, dest interface{}) error {
	val, ok := (*s.DB)[key]
	if !ok {
		return errors.New("no value")
	}
	deepCopy(val, dest)

	return nil
}

func (s StubDatabase) Set(key string, value interface{}) error {
	(*s.DB)[key] = value
	return nil
}

func (s StubDatabase) Incr(key string) (int64, error) {
	val, ok := (*s.DB)[key]
	if !ok {
		(*s.DB)[key] = int64(0)
		return 0, nil
	}
	val2, ok2 := val.(int64)
	if !ok2 {
		return 0, errors.New("value not int")
	}
	(*s.DB)[key] = val2 + 1
	return val2 + 1, nil
}
