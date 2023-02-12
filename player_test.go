package main

import (
	"testing"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		mem := NewStubDatabase()
		db := Memory{mem}
		var p Player
		db.GetPlayer(123, &p)

		p.NewGame(db)
		db.GetPlayer(p.ID, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, "11")
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		var p Player
		db.GetPlayer(123, &p)

		p.NewGame(db)

		db.GetPlayer(123, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, "02")

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p.NewGame(db)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)

		want = [3][3]Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		db.Set("gameID", 10)

		var p Player
		db.GetPlayer(123, &p)

		p.NewGame(db)
		db.GetPlayer(123, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)
		AssertString(t, game.ID, "10")

		p.NewGame(db)
		db.GetPlayer(123, &p)
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		AssertString(t, game.ID, "11")
	})
}
