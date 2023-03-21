package player

import (
	"reflect"
	"testing"

	"github.com/tumypmyp/chess/game"
)

func AssertGame(t testing.TB, got, want game.Game) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func AssertExactError(t testing.TB, a, b error) {
	t.Helper()
	if a != b {
		t.Fatalf("got error %v, but wanted: %v", a, b)
	}
}

func AssertBoard(t testing.TB, got, want [3][3]game.Mark) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}
