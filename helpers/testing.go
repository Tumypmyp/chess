package helpers

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

func AssertPlayerID(t testing.TB, got, want PlayerID) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}


func AssertExactError(t testing.TB, a, b error) {
	t.Helper()
	if a != b {
		t.Fatalf("got error %v, but wanted: %v", a, b)
	}
}