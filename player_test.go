package main

import (
	"reflect"
	"testing"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		t.Log("test start")
		p := NewPlayer(NewStubDatabase(), 123)

		t.Log(p)
		p.NewGame()
		t.Log(p)
		if p.CurrentGame() == nil {
			t.Fatalf("got no current game")
		}
		p.CurrentGame().Move("11")
		got := p.CurrentGame().Board
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
	t.Run("new game", func(t *testing.T) {

		p := NewPlayer(NewStubDatabase(), 123)
		p.NewGame()
		p.CurrentGame().Move("02")
		got := p.CurrentGame().Board
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
		p.NewGame()

		got = p.CurrentGame().Board
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
		AssertString(t, p.CurrentGame().ID, "10")

		p.NewGame()
		AssertString(t, p.CurrentGame().ID, "11")
	})
}

func AssertString(t testing.TB, got, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
