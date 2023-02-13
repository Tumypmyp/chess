package main

import (
	"testing"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		mem := NewStubDatabase()
		db := Memory{mem}
		p := NewPlayer(db, 123, "pl")

		p.NewGame(db, nil)
		db.GetPlayer(p.ID, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, "11", nil)
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		p := NewPlayer(db, 123, "pl")

		p.NewGame(db, nil)

		db.GetPlayer(123, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, "02", nil)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p.NewGame(db, nil)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)

		want = [3][3]Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		db.Set("gameID", int64(9))

		p := NewPlayer(db, 123, "pl")

		p.NewGame(db, nil)
		db.GetPlayer(123, &p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)
		AssertString(t, game.ID, "10")

		p.NewGame(db, nil)
		db.GetPlayer(123, &p)
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		AssertString(t, game.ID, "11")
	})
	t.Run(".NewGame updates player", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		id := int64(123456)
		p := NewPlayer(db, id, "pl")

		p.NewGame(db, nil)

		if len(p.GamesID) != 1 {
			t.Errorf("wanted 1 game, got %v", p.GamesID)
		}
	})

	t.Run("current game updates player", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		id := int64(123456)
		p := NewPlayer(db, id, "pl")

		NewGame(db, "123", nil, id)
		_, err := p.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("do", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		id := int64(1234)
		p := NewPlayer(db, 123, "pl")

		var err error
		err = p.Do(db, nil, "/new_game")
		AssertNoError(t, err)

		db.GetPlayer(id, &p)
		_, err = p.CurrentGame(db)
		AssertNoError(t, err)

		err = p.Do(db, nil, "/123")
		AssertError(t, err)
	})
	t.Run("do start game with other", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		p1 := NewPlayer(db, 123, "abc")
		p2 := NewPlayer(db, 456, "def")

		var err error
		err = p1.Do(db, nil, "/new_game @"+p2.Username)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("start game with other", func(t *testing.T) {
		db := Memory{NewStubDatabase()}
		p1 := NewPlayer(db, 123, "abc")
		p2 := NewPlayer(db, 456, "def")

		p1.NewGame(db, nil, p2.ID)

		var err error
		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
}
