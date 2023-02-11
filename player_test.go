package main

import (
	"testing"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		t.Log("test start")
		p := NewPlayer(NewStubDatabase(), 123)

		p.NewGame()
		game, err := p.CurrentGame()
		AssertNoError(t, err)

		p.Move("11")
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		p := NewPlayer(NewStubDatabase(), 123)

		p.NewGame()
		game, err := p.CurrentGame()
		AssertNoError(t, err)

		p.Move("02")
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p.NewGame()

		game, err = p.CurrentGame()
		AssertNoError(t, err)

		want = [3][3]Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
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
