package main

import (
	"encoding/json"
	"testing"
)

func TestStubDatabase(t *testing.T) {
	memory := NewStubDatabase()
	key := "abcd"
	value := "val"
	memory.Set(key, value)

	var got string
	memory.Get(key, &got)
	AssertString(t, got, value)

}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
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
	deepCopy((*s.DB)[key], dest)
	return nil
}
func (s StubDatabase) Set(key string, value interface{}) error {
	(*s.DB)[key] = value
	return nil
}
