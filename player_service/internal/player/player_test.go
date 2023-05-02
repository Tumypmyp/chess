package player

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
)

func TestPlayer(t *testing.T) {
	t.Run("new player", func(t *testing.T) {
		id := int64(123)
		p := NewPlayer(id, "pl")
		want := Player{
			ID: id,
			Username: "pl",
		}
		AssertPlayer(t, p, want)
	})

	// t.Run("move player", func(t *testing.T) {
	// 	mem := memory.NewStubDatabase()
	// 	db := mem
	// 	p := NewPlayer(db, 123, "pl")

	// 	t.Log("player created")
	// 	err := AddGameToPlayer(db, p.ID)
	// 	AssertNoError(t, err)

	// 	t.Log(mem.DB)
	// 	t.Log(p)

	// 	_, err = CurrentGame(p.ID, db)
	// 	AssertNoError(t, err)

	// 	Do(p.ID, db, "11", 0)
	// 	_, err = CurrentGame(p.ID, db)
	// 	AssertNoError(t, err)
	// 	// want := [3][3]g.Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
	// 	// AssertBoard(t, game.Board, want)
	// })
	t.Run("new game", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(123)
		_ = MakePlayer(db, id, "pl")

		AddGameToPlayer(db, 0, id)

		_, err := getPlayer(id, db)
		AssertNoError(t, err)

		_, err = CurrentGame(id, db)
		AssertNoError(t, err)

		Do(id, db, "02", 0)

		_, err = CurrentGame(id, db)
		AssertNoError(t, err)
		// want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		// AssertBoard(t, game.Board, want)

		AddGameToPlayer(db, 1, id)

		_, err = CurrentGame(id, db)
		AssertNoError(t, err)

		// want = [3][3]g.Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		// AssertBoard(t, game.Board, want)
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := memory.NewStubDatabase()
		db.Set("gameID", int64(9))

		id := int64(123)
		_ = MakePlayer(db, id, "pl")

		AddGameToPlayer(db, 10, id)
		_, err := getPlayer(id, db)
		AssertNoError(t, err)

		game, err := CurrentGame(id, db)
		AssertNoError(t, err)
		AssertInt(t, game, 10)

		AddGameToPlayer(db, 11, id)
		_, err = getPlayer(id, db)
		AssertNoError(t, err)

		game, err = CurrentGame(id, db)
		AssertNoError(t, err)
		AssertInt(t, game, 11)
	})
	t.Run(".NewGame do not updates player", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(123456)
		p := MakePlayer(db, id, "pl")

		AddGameToPlayer(db, 1, p.ID)

		if len(p.GamesID) != 0 {
			t.Errorf("wanted 0 game, got %v", p.GamesID)
		}
		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		if len(p.GamesID) != 1 {
			t.Errorf("wanted 1 game, got %v", p.GamesID)
		}

	})

	t.Run("current game ", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(123456)
		_ = MakePlayer(db, id, "pl")

		AddGameToPlayer(db, 1, id)
		_, err := CurrentGame(id, db)
		AssertNoError(t, err)

	})
}

// func TestCmd(t *testing.T) {
// 	t.Run("2 players play in turns", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p1 := NewPlayer(db, int64(12), "pl12")
// 		p2 := NewPlayer(db, int64(13), "pl13")
// 		AddGameToPlayer(db, 1, p1.ID, p2.ID)

// 		var err error
// 		_, err = Do(p2.ID, db, "11", 0)
// 		AssertError(t, err)

// 		_, err = Do(p1.ID, db, "11", 0)
// 		AssertNoError(t, err)
// 	})

// 	t.Run("do", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		id := int64(123)
// 		_ = NewPlayer(db, id, "pl")

// 		var err error
// 		cmd := "newgame"
// 		_, err = Cmd(db, cmd, "/"+cmd, id, 0)

// 		// _, err = Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
// 		// 	{Type: "bot_command", Offset: 0, Length: len(cmd)},
// 		// }}, p, 0)
// 		AssertNoError(t, err)

// 		_, err = getPlayer(id, db)
// 		AssertNoError(t, err)

// 		_, err = getPlayer(int64(456), db)
// 		AssertError(t, err)

// 		_, err = CurrentGame(id, db)
// 		AssertNoError(t, err)

// 		_, err = Do(id, db, "/123", 0)
// 		AssertError(t, err)
// 	})
// 	t.Run("do move", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		id := int64(12)
// 		_ = NewPlayer(db, id, "pl")
// 		_ = AddGameToPlayer(db, 1, id)
// 		_, err := getPlayer(id, db)
// 		AssertNoError(t, err)

// 		_, err = Do(id, db, "00", 0)
// 		AssertNoError(t, err)

// 		_, err = CurrentGame(id, db)
// 		AssertNoError(t, err)

// 		// board := [3][3]g.Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
// 		// want := g.Game{PlayersID: []int64{id},
// 		// 	Description:   "@pl ",
// 		// 	CurrentPlayer: 0,
// 		// 	ChatsID:       []int64{12},
// 		// 	Board:         board,
// 		// 	ID:            0}
// 		// AssertGame(t, game, want)
// 	})
// 	t.Run("do start game with other", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p1 := NewPlayer(db, int64(123), "abc")
// 		p2 := NewPlayer(db, int64(456), "def")

// 		var err error
// 		cmd := "newgame"
// 		_, err = Cmd(db, cmd, "/"+cmd+" @"+p2.Username, p1.ID, 0)
// 		// _, err = Cmd(db, &tgbotapi.Message{Text: cmd + " @" + p2.Username, Entities: []tgbotapi.MessageEntity{
// 		// 	{Type: "bot_command", Offset: 0, Length: len(cmd)},
// 		// }}, p1, 0)
// 		AssertNoError(t, err)

// 		_, err = CurrentGame(p1.ID, db)
// 		AssertNoError(t, err)
// 		_, err = CurrentGame(p2.ID, db)
// 		AssertNoError(t, err)

// 	})
// 	t.Run("start game with other", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p1 := NewPlayer(db, int64(123), "abc")
// 		p2 := NewPlayer(db, int64(456), "def")

// 		AddGameToPlayer(db, 1, p1.ID, p2.ID)

// 		var err error
// 		_, err = CurrentGame(p1.ID, db)
// 		AssertNoError(t, err)
// 		_, err = CurrentGame(p2.ID, db)
// 		AssertNoError(t, err)

// 	})
// 	t.Run("start game with 2 others", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p1 := NewPlayer(db, int64(123), "abc")
// 		p2 := NewPlayer(db, int64(456), "def")
// 		p3 := NewPlayer(db, int64(789), "ghi")

// 		var err error
// 		cmd := "newgame"
// 		_, err = Cmd(db, cmd, "/"+cmd+" @"+p2.Username+" @"+p3.Username, p1.ID, 0)
// 		// _, err = Cmd(db, &tgbotapi.Message{Text: cmd + " @" + p2.Username + " @" + p3.Username,
// 		// 	Entities: []tgbotapi.MessageEntity{
// 		// 		{Type: "bot_command", Offset: 0, Length: len(cmd)},
// 		// 	}}, p1, 0)
// 		AssertNoError(t, err)

// 		_, err = CurrentGame(p1.ID, db)
// 		AssertNoError(t, err)
// 		_, err = CurrentGame(p2.ID, db)
// 		AssertNoError(t, err)
// 		_, err = CurrentGame(p3.ID, db)
// 		AssertNoError(t, err)

// 	})

// 	t.Run("other player with different chat", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p := NewPlayer(db, int64(123), "pl")

// 		AddGameToPlayer(db, 1, p.ID)

// 		p, err := getPlayer(p.ID, db)
// 		AssertNoError(t, err)

// 		_, err = CurrentGame(p.ID, db)
// 		AssertNoError(t, err)

// 		Do(p.ID, db, "02", 0)

// 		_, err = CurrentGame(p.ID, db)
// 		AssertNoError(t, err)
// 		// want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
// 		// AssertBoard(t, game.Board, want)

// 		p2 := NewPlayer(db, int64(456), "pl")

// 		_, err = CurrentGame(p2.ID, db)
// 		AssertExactError(t, err, NoCurrentGameError{})

// 		AddGameToPlayer(db, 2, p2.ID)
// 		_, err = CurrentGame(p2.ID, db)

// 		AssertNoError(t, err)
// 		// want = [3][3]g.Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
// 		// AssertBoard(t, game.Board, want)

// 		_, err = CurrentGame(p.ID, db)
// 		AssertNoError(t, err)
// 		// want = [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
// 		// AssertBoard(t, game.Board, want)
// 	})
// }

// func TestPlayerResponses(t *testing.T) {
// 	t.Run("start player", func(t *testing.T) {
// 		db := memory.NewStubDatabase()

// 		p := NewPlayer(db, int64(123), "pl")
// 		r, err := Do(p.ID, db, "12", 0)
// 		AssertError(t, err)
// 		AssertString(t, r.Text, NoCurrentGameError{}.Error())
// 	})
// 	t.Run("start game with other", func(t *testing.T) {
// 		db := memory.NewStubDatabase()

// 		p1 := NewPlayer(db, int64(123), "abc")
// 		p2 := NewPlayer(db, int64(456), "def")

// 		err := AddGameToPlayer(db, 1, p1.ID, p2.ID)
// 		AssertNoError(t, err)
// 		// AssertInt(t, int64(len(r.ChatsID)), 2)

// 		r, err := Do(p1.ID, db, "11", 0)
// 		AssertNoError(t, err)
// 		// AssertString(t, r2.Text, "Started")
// 		AssertInt(t, int64(len(r.ChatsID)), 2)
// 	})
// 	t.Run("start game with 2 other", func(t *testing.T) {
// 		db := memory.NewStubDatabase()

// 		p1 := NewPlayer(db, int64(123), "abc")
// 		p2 := NewPlayer(db, int64(456), "def")
// 		p3 := NewPlayer(db, int64(789), "ghi")

// 		err := AddGameToPlayer(db, 1, p1.ID, p2.ID, p3.ID)
// 		AssertNoError(t, err)
// 		// AssertInt(t, int64(len(r.ChatsID)), 3)

// 		r, err := Do(p1.ID, db, "11", 0)
// 		AssertNoError(t, err)
// 		AssertInt(t, int64(len(r.ChatsID)), 3)
// 	})
// }

// func TestPlayerCmd(t *testing.T) {
// 	t.Run("newgame", func(t *testing.T) {
// 		db := memory.NewStubDatabase()

// 		p1 := NewPlayer(db, int64(123), "abc")

// 		var err error
// 		cmd := "newgame"
// 		r, err := Cmd(db, cmd, "/"+cmd, p1.ID, 0)
// 		// r, err := Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
// 		// 	{Type: "bot_command", Offset: 0, Length: len(cmd)},
// 		// }}, p1, 0)
// 		AssertNoError(t, err)
// 		_, err = CurrentGame(p1.ID, db)
// 		AssertNoError(t, err)
// 		AssertInt(t, int64(len(r.ChatsID)), 1)

// 	})
// 	t.Run("leaderboard", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		//
// 		p1 := NewPlayer(db, int64(123), "abc")

// 		var err error
// 		cmd := "leaderboard"
// 		r, err := Cmd(db, cmd, "/"+cmd, p1.ID, 0)
// 		// r, err := Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
// 		// 	{Type: "bot_command", Offset: 0, Length: len(cmd)},
// 		// }}, p1, 0)
// 		AssertExactError(t, err, NoConnectionError{})
// 		AssertString(t, r.Text, NoConnectionError{}.Error())

// 		_, err = CurrentGame(p1.ID, db)
// 		AssertExactError(t, err, NoCurrentGameError{})

// 	})
// 	t.Run("no such command", func(t *testing.T) {
// 		db := memory.NewStubDatabase()

// 		p1 := NewPlayer(db, int64(123), "abc")

// 		cmd1 := "command"
// 		r, err := Cmd(db, cmd1, "/"+cmd1, p1.ID, 0)

// 		AssertExactError(t, err, NoSuchCommandError{"command"})
// 		AssertString(t, r.Text, NoSuchCommandError{"command"}.Error())

// 		cmd2 := "newgame2"
// 		r, err = Cmd(db, cmd2, "/"+cmd2, p1.ID, 0)

// 		AssertExactError(t, err, NoSuchCommandError{"newgame2"})
// 		AssertInt(t, int64(len(r.ChatsID)), 1)
// 		AssertString(t, r.Text, NoSuchCommandError{"newgame2"}.Error())

// 		_, err = CurrentGame(p1.ID, db)
// 		AssertExactError(t, err, NoCurrentGameError{})
// 	})

// 	t.Run("player do new game", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		p1 := NewPlayer(db, int64(123), "abc")

// 		err := AddGameToPlayer(db, 0, p1.ID)
// 		AssertNoError(t, err)
// 		// AssertString(t, r.Text, "@abc \nStarted\n")
// 	})
// 	t.Run("cmd to players id", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		_ = NewPlayer(db, int64(123), "abc")
// 		_ = NewPlayer(db, int64(456), "def")

// 		players, err := cmdToPlayersID(db, "/newgame @abc @def")
// 		want := []int64{123, 456}
// 		AssertNoError(t, err)
// 		if !reflect.DeepEqual(players, want) {
// 			t.Errorf("got %v, wanted %v", players, want)
// 		}
// 	})
// 	t.Run("cmd to no players id", func(t *testing.T) {
// 		db := memory.NewStubDatabase()
// 		_ = NewPlayer(db, int64(123), "abc")
// 		_ = NewPlayer(db, int64(456), "def")

// 		_, err := cmdToPlayersID(db, "/newgame @aaa @bbb")
// 		AssertExactError(t, err, NoUsernameInDatabaseError{})
// 	})
// }
