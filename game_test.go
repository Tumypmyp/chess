package main

import (
	"testing"
)

func TestGame(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		db := NewStubDatabase()
		player := NewPlayer(db, PlayerID{12,12}, "pl")
		game := NewGame(db, nil, player.ID)
		err := game.Move(player, "00")
		AssertNoError(t, err)

		player.Get(player.ID, db)
		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{player.ID},
			Description:   "@pl ",
			CurrentPlayer: 0,
			Status:        Started,
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)
	})
	t.Run("do move", func(t *testing.T) {
		db := NewStubDatabase()
		player := NewPlayer(db, PlayerID{12,12}, "pl")
		game := NewGame(db, nil, player.ID)
		err := player.Do(db, nil, "00")
		AssertNoError(t, err)

		game, err = player.CurrentGame(db)
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{player.ID},
			Description:   "@pl ",
			CurrentPlayer: 0,
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)
	})
	t.Run("2 players play in turns", func(t *testing.T) {
		db := NewStubDatabase()
		p1 := NewPlayer(db, PlayerID{12,12}, "pl12")
		p2 := NewPlayer(db, PlayerID{13,13}, "pl13")
		p1.NewGame(db, nil, p2.ID)

		var err error
		err = p2.Move(db, nil, "11")
		AssertError(t, err)

		err = p1.Move(db, nil, "11")
		AssertNoError(t, err)

	})
	t.Run("2 players", func(t *testing.T) {
		mem := NewStubDatabase()
		db := mem
		p1 := NewPlayer(db, PlayerID{12,12}, "pl12")
		p2 := NewPlayer(db, PlayerID{13,13}, "pl13")
		p1.NewGame(db, nil, p2.ID)

		p1.Get(p1.ID, db)
		p2.Get(p2.ID, db)
		game, err := p1.CurrentGame(db)
		AssertNoError(t, err)

		err = game.Move(p1, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{p1.ID, p2.ID},
			Description:   "@pl12 @pl13 ",
			CurrentPlayer: 1,
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)

		err = game.Move(p2, "01")
		AssertNoError(t, err)

		board = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		want.Board = board
		want.CurrentPlayer = 0
		AssertGame(t, game, want)
		got := game.String()
		str := "@pl12 @pl13 \nStarted\nXO-\n---\n---\n"
		AssertString(t, got, str)

	})

	t.Run("game status", func(t *testing.T) {
		db := NewStubDatabase()
		player := NewPlayer(db, PlayerID{12,12}, "pl")
		game := NewGame(db, nil, player.ID)
		AssertStatus(t, game.Status, Started)

		err := game.Move(player, "00")
		AssertNoError(t, err)

		err = game.Move(player, "11")
		AssertNoError(t, err)
		err = game.Move(player, "22")
		AssertNoError(t, err)
		AssertStatus(t, game.Status, Finished)
	})
}
