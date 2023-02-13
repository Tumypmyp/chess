package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		var player Player
		db.GetPlayer(12, &player)
		game := NewGame(db, "122", nil, player.ID)
		game.Move(player.ID, "00")

		want := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("2 players", func(t *testing.T) {
		mem := NewStubDatabase()
		db := Memory{mem}
		var p1 Player
		var p2 Player
		db.GetPlayer(12, &p1)
		db.GetPlayer(13, &p2)
		p1.NewGame(db, nil, p2.ID)

		db.GetPlayer(p1.ID, &p1)
		db.GetPlayer(p2.ID, &p2)
		game, err := p1.CurrentGame(db)
		AssertNoError(t, err)

		err = game.Move(p1.ID, "00")
		AssertNoError(t, err)

		want := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		err = game.Move(p2.ID, "01")
		AssertNoError(t, err)

		want = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

	})
}
