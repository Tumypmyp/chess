package main

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"testing"
)

func TestMemory(t *testing.T) {

	t.Run("get", func(t *testing.T) {
		db := NewStubDatabase()
		var player Player
		err := player.Get(12, db)
		AssertError(t, err)
		player = NewPlayer(db, 12, "abc")

		t.Log(player)
		if player.ID != 12 {
			t.Errorf("got %v, want %v", player.ID, 12)
		}
	})
	t.Run("set/get", func(t *testing.T) {
		mem := NewStubDatabase()
		p := Player{ID: 1234}
		p.Store(mem)
		var got Player
		got.Get(1234, mem)
		if !reflect.DeepEqual(p, got) {
			t.Errorf("got %v, want %v", got, p)
		}
	})
}
func TestStubDatabase(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		memory := NewStubDatabase()
		key := "abcd"
		value := "val"
		memory.Set(key, value)

		var got string
		memory.Get(key, &got)
		AssertString(t, got, value)
	})
	t.Run("player", func(t *testing.T) {
		memory := NewStubDatabase()
		key := "abcd"
		value := Player{ID: 1234}
		memory.Set(key, value)
		t.Logf("mem: %v", *memory.DB)

		var got Player
		memory.Get(key, &got)
		if !reflect.DeepEqual(value, got) {
			t.Errorf("got %v, want %v", got, value)
		}
	})
	t.Run("player with game", func(t *testing.T) {
		memory := NewStubDatabase()
		key := "abcd"
		value := Player{ID: 1234, GamesID: []int64{12}}

		memory.Set(key, value)

		var got Player
		memory.Get(key, &got)
		if !reflect.DeepEqual(value, got) {
			t.Errorf("got %v, want %v", got, value)
		}
	})
}

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
		return 0, errors.New("no value")
	}
	val2, ok2 := val.(int64)
	if !ok2 {
		return 0, errors.New("value not int")
	}
	(*s.DB)[key] = val2 + 1
	return val2 + 1, nil
}
