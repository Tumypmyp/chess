package player

import (
	"reflect"
	"testing"

	"github.com/tumypmyp/chess/player_service/internal/game"
)

func AssertGame(t testing.TB, got, want game.Game) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}


func AssertBoard(t testing.TB, got, want [3][3]game.Mark) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}


func AssertPlayer(t testing.TB, got, want Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
