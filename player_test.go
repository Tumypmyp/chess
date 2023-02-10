package main

import (
	"reflect"
	"testing"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		t.Log("test start")
		p := NewPlayer(NewStubDatabase(), 123)

		p.NewGame()
		game, err := p.CurrentGame()
		AssertNoError(t, err)

		game.Move("11")
		got := game.Board
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("new game", func(t *testing.T) {
		p := NewPlayer(NewStubDatabase(), 123)

		p.NewGame()
		game, err := p.CurrentGame()
		AssertNoError(t, err)

		game.Move("02")
		got := game.Board
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}

		p.NewGame()

		game, err = p.CurrentGame()
		AssertNoError(t, err)

		got = game.Board
		want = [3][3]Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := NewStubDatabase()
		db.Set("gameID", 10)

		p := NewPlayer(db, 123)

		p.NewGame()
		game, err := p.CurrentGame()
		AssertNoError(t, err)
		AssertString(t, game.ID, "10")

		p.NewGame()
		game, err = p.CurrentGame()
		AssertNoError(t, err)
		AssertString(t, game.ID, "11")
	})
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
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
