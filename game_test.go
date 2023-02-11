package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		player := NewPlayer(nil, 12)
		game := NewGame(12, &player)
		game.Move(&player, "00")

		want := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("2 players", func(t *testing.T) {
		db := NewStubDatabase()
		p1 := NewPlayer(db, 12)
		p2 := NewPlayer(db, 13)
		p1.NewGame(&p2)

		game, err := p1.CurrentGame()
		AssertNoError(t, err)

		err = game.Move(&p1, "00")
		AssertNoError(t, err)

		want := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		err = game.Move(&p2, "01")
		AssertNoError(t, err)

		want = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

	})
}
