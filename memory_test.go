package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	t.Run("set and ", func(t *testing.T) {
		/*
			t.Log("base adding")
			memory := NewDatabase()

			t.Log("base added")
			err := memory.Set("test", "value")
			if err != nil {
				t.Fatalf("cant set %v", err)
			}
			var got string
			err = memory.Get("test", &got)
			if err != nil {
				t.Fatalf("cant get %v", err)
			}
			want := "value"
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %v, want %v", got, want)
			}
		*/
	})

}

func TestStubDatabase(t *testing.T) {
	memory := NewStubDatabase()
	key := "abcd"
	memory.Set(key, "value")
	t.Log(memory)
	var got string
	memory.Get(key, &got)
	want := "value"
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

}

func deepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}

type StubDatabase map[string]interface{}

func NewStubDatabase() StubDatabase {
	return make(map[string]interface{})
}

func (s StubDatabase) Get(key string, dest interface{}) error {
	deepCopy(s[key], dest)
	return nil
}
func (s StubDatabase) Set(key string, value interface{}) error {
	s[key] = value
	return nil
}
