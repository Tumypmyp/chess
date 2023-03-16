package game

import (
	"reflect"
	"testing"
	
)

func AssertInt(t testing.TB, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

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

func AssertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error but did not get one: %v", err)
	}
}

func AssertBoard(t testing.TB, got, want [3][3]Mark) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
