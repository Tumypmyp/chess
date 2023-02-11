package main

import (
	"reflect"
	"testing"
)

func AssertString(t testing.TB, got, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one: %v", err)
	}
}
func AssertBoard(t testing.TB, got, want [3][3]Mark) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
