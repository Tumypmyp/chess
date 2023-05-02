package player

import (
	"reflect"
	"testing"
)

func AssertPlayer(t testing.TB, got, want Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
