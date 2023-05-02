package game

import (
	"reflect"
	"testing"
)

func AssertStatus(t testing.TB, got, want GameStatus) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertGame(t testing.TB, got, want Game) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func AssertBoard(t testing.TB, got, want [3][3]Mark) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
