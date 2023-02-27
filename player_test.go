package main

import (
	"testing"
	"fmt"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		mem := NewStubDatabase()
		db := mem
		p := NewPlayer(db, PlayerID{123,123}, "pl")

		p.NewGame(db, nil)
		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "11")
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		db := NewStubDatabase()
		p := NewPlayer(db, PlayerID{123,123}, "pl")

		p.NewGame(db, nil)

		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "02")

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
		db := NewStubDatabase()
		db.Set("gameID", int64(9))

		p := NewPlayer(db, PlayerID{123,123}, "pl")

		p.NewGame(db, nil)
		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)
		AssertInt(t, game.ID, 10)

		p.NewGame(db, nil)
		p.Get(p.ID, db)
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		AssertInt(t, game.ID, 11)
	})
	t.Run(".NewGame updates player", func(t *testing.T) {
		db := NewStubDatabase()
		id := PlayerID{123456, 123456}
		p := NewPlayer(db, id, "pl", )

		p.NewGame(db, nil)

		if len(p.GamesID) != 1 {
			t.Errorf("wanted 1 game, got %v", p.GamesID)
		}
	})

	t.Run("current game updates player", func(t *testing.T) {
		db := NewStubDatabase()
		id := PlayerID{123456, 123456}
		p := NewPlayer(db, id, "pl")

		NewGame(db, nil, p.ID)
		_, err := p.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("do", func(t *testing.T) {
		db := NewStubDatabase()
		id := PlayerID{1234, 1234}
		//?? why 1234-123
		p := NewPlayer(db, PlayerID{123, 123}, "pl")

		var err error
		err = p.Do(db, nil, "/new_game")
		AssertNoError(t, err)

		p.Get(id, db)
		_, err = p.CurrentGame(db)
		AssertNoError(t, err)

		err = p.Do(db, nil, "/123")
		AssertError(t, err)
	})
	t.Run("do start game with other", func(t *testing.T) {
		db := NewStubDatabase()
		p1 := NewPlayer(db, PlayerID{123, 123}, "abc")
		p2 := NewPlayer(db, PlayerID{456, 456}, "def")

		var err error
		err = p1.Do(db, nil, "/new_game @"+p2.Username)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("start game with other", func(t *testing.T) {
		db := NewStubDatabase()
		p1 := NewPlayer(db, PlayerID{123, 123}, "abc")
		p2 := NewPlayer(db, PlayerID{456, 456}, "def")

		p1.NewGame(db, nil, p2.ID)

		var err error
		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("start game with 2 others", func(t *testing.T) {
		db := NewStubDatabase()
		p1 := NewPlayer(db, PlayerID{123, 123}, "abc")
		p2 := NewPlayer(db, PlayerID{456, 456}, "def")
		p3 := NewPlayer(db, PlayerID{789, 789}, "ghi")
		

		var err error
		err = p1.Do(db, nil, "/new_game @"+p2.Username+" @"+p3.Username)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p3.CurrentGame(db)
		AssertNoError(t, err)

	})

	t.Run("one player with different chat", func(t *testing.T) {
		db := NewStubDatabase()
		p := NewPlayer(db, PlayerID{123, 123}, "pl")

		p.NewGame(db, nil)

		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "02")

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		
		t.Log(db.DB)
		
		p2 := NewPlayer(db, PlayerID{12345, 123}, "pl")
		fmt.Println(p2)
		_, err = p2.CurrentGame(db)
		
		fmt.Println(p2)
		AssertError(t, err)

		p2.NewGame(db, nil)

		t.Log(db.DB)
		fmt.Println(p2)
		game, err = p2.CurrentGame(db)
		fmt.Println(p2)
		AssertNoError(t, err)
		want = [3][3]Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		
		game, err = p.CurrentGame(db)
		fmt.Println(p)
		AssertNoError(t, err)
		want = [3][3]Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
}
