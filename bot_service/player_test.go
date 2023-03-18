package main

import (
	"testing"

	g "github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/memory"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		p := NewPlayer(db, g.PlayerID{123}, "pl")

		p.NewGame(db, nil)
		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "11")
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p := NewPlayer(db, g.PlayerID{123}, "pl")

		p.NewGame(db, nil)

		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "02")

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p.NewGame(db, nil)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)

		want = [3][3]g.Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := memory.NewStubDatabase()
		db.Set("gameID", int64(9))

		p := NewPlayer(db, g.PlayerID{123}, "pl")

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
		db := memory.NewStubDatabase()
		id := g.PlayerID{123456}
		p := NewPlayer(db, id, "pl")

		p.NewGame(db, nil)

		if len(p.GamesID) != 1 {
			t.Errorf("wanted 1 game, got %v", p.GamesID)
		}
	})

	t.Run("current game updates player", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		id := g.PlayerID{123456}
		p := NewPlayer(db, id, "pl")

		p.NewGame(db, nil, id)
		t.Log(mem.DB)
		_, err := p.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("do", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := g.PlayerID{1234}
		//?? why 1234-123
		p := NewPlayer(db, g.PlayerID{123}, "pl")

		var err error
		err = p.Do(db, nil, "/newgame")
		AssertNoError(t, err)

		p.Get(id, db)
		_, err = p.CurrentGame(db)
		AssertNoError(t, err)

		err = p.Do(db, nil, "/123")
		AssertError(t, err)
	})
	t.Run("do start game with other", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		p1 := NewPlayer(db, g.PlayerID{123}, "abc")
		p2 := NewPlayer(db, g.PlayerID{456}, "def")

		var err error
		err = p1.Do(db, nil, "/newgame @"+p2.Username)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		t.Log(mem.DB)
		AssertNoError(t, err)

	})
	t.Run("start game with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, g.PlayerID{123}, "abc")
		p2 := NewPlayer(db, g.PlayerID{456}, "def")

		p1.NewGame(db, nil, p2.ID)

		var err error
		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("start game with 2 others", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, g.PlayerID{123}, "abc")
		p2 := NewPlayer(db, g.PlayerID{456}, "def")
		p3 := NewPlayer(db, g.PlayerID{789}, "ghi")

		var err error
		err = p1.Do(db, nil, "/newgame @"+p2.Username+" @"+p3.Username)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p3.CurrentGame(db)
		AssertNoError(t, err)

	})

	t.Run("other player with different chat", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p := NewPlayer(db, g.PlayerID{123}, "pl")

		p.NewGame(db, nil)

		p.Get(p.ID, db)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Move(db, nil, "02")

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p2 := NewPlayer(db, g.PlayerID{456}, "pl")

		_, err = p2.CurrentGame(db)
		AssertError(t, err)

		p2.NewGame(db, nil)
		game, err = p2.CurrentGame(db)

		AssertNoError(t, err)
		want = [3][3]g.Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want = [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
}

func TestPlayerResponses(t *testing.T) {
	t.Run("start player", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := NewStubBot()

		p := NewPlayer(db, g.PlayerID{123}, "pl")
		err := p.Do(db, bot, "12")
		AssertError(t, err)
		AssertString(t, bot.Read(), NoCurrentGameError{}.Error())
	})
	t.Run("start game with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := NewStubBot()
		p1 := NewPlayer(db, g.PlayerID{123}, "abc")
		p2 := NewPlayer(db, g.PlayerID{456}, "def")

		p1.NewGame(db, bot, p2.ID)
		AssertInt(t, bot.Len(), 2)

		p1.Do(db, bot, "11")
		AssertInt(t, bot.Len(), 4)
	})
	t.Run("start game with 2 other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := NewStubBot()
		p1 := NewPlayer(db, g.PlayerID{123}, "abc")
		p2 := NewPlayer(db, g.PlayerID{456}, "def")
		p3 := NewPlayer(db, g.PlayerID{789}, "ghi")

		p1.NewGame(db, bot, p2.ID, p3.ID)
		AssertInt(t, bot.Len(), 3)

		p1.Do(db, bot, "11")
		AssertInt(t, bot.Len(), 6)
	})
}
