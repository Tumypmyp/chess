package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		player := NewPlayer(db, 12, "pl")
		game := NewGame(db, "122", nil, player.ID)
		game.Move(player.ID, "00")

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []int64{12},
			PlayersUsername: []string{"pl"},
			Board:           board,
			ID:              "122"}
		AssertGame(t, game, want)
	})
	t.Run("2 players", func(t *testing.T) {
		mem := NewStubDatabase()
		db := Memory{mem}
		p1 := NewPlayer(db, 12, "pl12")
		p2 := NewPlayer(db, 13, "pl13")
		p1.NewGame(db, nil, p2.ID)

		db.GetPlayer(p1.ID, &p1)
		db.GetPlayer(p2.ID, &p2)
		game, err := p1.CurrentGame(db)
		AssertNoError(t, err)

		err = game.Move(p1.ID, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []int64{12, 13},
			PlayersUsername: []string{"pl12", "pl13"},
			Board:           board,
			ID:              "0"}
		AssertGame(t, game, want)

		err = game.Move(p2.ID, "01")
		AssertNoError(t, err)

		board = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		want.Board = board
		AssertGame(t, game, want)
		got := game.String()
		str := "@pl12 @pl13 \nXO-\n---\n---\n"
		AssertString(t, got, str)

	})
}
