package player

import (
	"reflect"
	"testing"

	g "github.com/tumypmyp/chess/game"
	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestPlayer(t *testing.T) {
	t.Run("move player", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		p := NewPlayer(db, PlayerID(123), "pl")

		NewGame(db, p.ID)

		t.Log(mem.DB)
		t.Log(p)

		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		t.Log(mem.DB)
		t.Log(p)
		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Do(db, "11", 0)
		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p := NewPlayer(db, PlayerID(123), "pl")

		NewGame(db, p.ID)

		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Do(db, "02", 0)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		NewGame(db, p.ID)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)

		want = [3][3]g.Mark{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)
	})
	t.Run("new game is next id", func(t *testing.T) {
		db := memory.NewStubDatabase()
		db.Set("gameID", int64(9))

		p := NewPlayer(db, PlayerID(123), "pl")

		NewGame(db, p.ID)
		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		game, err := p.CurrentGame(db)
		AssertNoError(t, err)
		AssertInt(t, game.ID, 10)

		NewGame(db, p.ID)
		p, err = getPlayer(p.ID, db)
		AssertNoError(t, err)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		AssertInt(t, game.ID, 11)
	})
	t.Run(".NewGame do not updates player", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := PlayerID(123456)
		p := NewPlayer(db, id, "pl")

		NewGame(db, p.ID)

		if len(p.GamesID) != 0 {
			t.Errorf("wanted 0 game, got %v", p.GamesID)
		}
		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		if len(p.GamesID) != 1 {
			t.Errorf("wanted 1 game, got %v", p.GamesID)
		}

	})

	t.Run("current game updates player", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		id := PlayerID(123456)
		p := NewPlayer(db, id, "pl")

		NewGame(db, id)
		t.Log(mem.DB)
		_, err := p.CurrentGame(db)
		AssertNoError(t, err)

	})

	t.Run("2 players play in turns", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, PlayerID(12), "pl12")
		p2 := NewPlayer(db, PlayerID(13), "pl13")
		NewGame(db, p1.ID, p2.ID)

		var err error
		_, err = p2.Do(db, "11", 0)
		AssertError(t, err)

		_, err = p1.Do(db, "11", 0)
		AssertNoError(t, err)
	})

	t.Run("do", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := PlayerID(123)
		p := NewPlayer(db, PlayerID(123), "pl")

		var err error
		cmd := "/newgame"
		_, err = Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}, p, 0)
		AssertNoError(t, err)

		p, err = getPlayer(id, db)
		AssertNoError(t, err)

		_, err = getPlayer(PlayerID(456), db)
		AssertError(t, err)

		t.Log(db.DB)
		// p, _ = getPlayer(p.ID, db)
		_, err = p.CurrentGame(db)
		AssertNoError(t, err)

		_, err = p.Do(db, "/123", 0)
		AssertError(t, err)
	})
	t.Run("do move", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := NewPlayer(db, PlayerID(12), "pl")
		_ = NewGame(db, player.ID)
		player, err := getPlayer(player.ID, db)
		AssertNoError(t, err)
		t.Log(db.DB)
		t.Log(player)
		_, err = player.Do(db, "00", 0)
		AssertNoError(t, err)

		game, err := player.CurrentGame(db)
		AssertNoError(t, err)

		board := [3][3]g.Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := g.Game{PlayersID: []PlayerID{player.ID},
			Description:   "@12 ",
			CurrentPlayer: 0,
			ChatsID:       []int64{12},
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)
	})
	t.Run("do start game with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, PlayerID(123), "abc")
		p2 := NewPlayer(db, PlayerID(456), "def")

		var err error
		cmd := "/newgame"
		_, err = Cmd(db, &tgbotapi.Message{Text: cmd + " @" + p2.Username, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}, p1, 0)
		AssertNoError(t, err)

		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		t.Log(db.DB)
		AssertNoError(t, err)

	})
	t.Run("start game with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, PlayerID(123), "abc")
		p2 := NewPlayer(db, PlayerID(456), "def")

		NewGame(db, p1.ID, p2.ID)

		var err error
		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		_, err = p2.CurrentGame(db)
		AssertNoError(t, err)

	})
	t.Run("start game with 2 others", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, PlayerID(123), "abc")
		p2 := NewPlayer(db, PlayerID(456), "def")
		p3 := NewPlayer(db, PlayerID(789), "ghi")

		var err error
		cmd := "/newgame"
		_, err = Cmd(db, &tgbotapi.Message{Text: cmd + " @" + p2.Username + " @" + p3.Username,
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: len(cmd)},
			}}, p1, 0)
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
		p := NewPlayer(db, PlayerID(123), "pl")

		NewGame(db, p.ID)

		p, err := getPlayer(p.ID, db)
		AssertNoError(t, err)

		game, err := p.CurrentGame(db)
		AssertNoError(t, err)

		p.Do(db, "02", 0)

		game, err = p.CurrentGame(db)
		AssertNoError(t, err)
		want := [3][3]g.Mark{{0, 0, 1}, {0, 0, 0}, {0, 0, 0}}
		AssertBoard(t, game.Board, want)

		p2 := NewPlayer(db, PlayerID(456), "pl")

		_, err = p2.CurrentGame(db)
		AssertExactError(t, err, NoCurrentGameError{})

		NewGame(db, p2.ID)
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

		p := NewPlayer(db, PlayerID(123), "pl")
		r, err := p.Do(db, "12", 0)
		AssertError(t, err)
		AssertString(t, r[0].Text, NoCurrentGameError{}.Error())
	})
	t.Run("start game with other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		
		p1 := NewPlayer(db, PlayerID(123), "abc")
		p2 := NewPlayer(db, PlayerID(456), "def")

		r := NewGame(db, p1.ID, p2.ID)
		AssertInt(t, int64(len(r)), 2)

		r,  err := p1.Do(db, "11", 0)
		AssertNoError(t, err)
		// AssertString(t, r2.Text, "Started")
		AssertInt(t, int64(len(r)), 2)
	})
	t.Run("start game with 2 other", func(t *testing.T) {
		db := memory.NewStubDatabase()
		
		p1 := NewPlayer(db, PlayerID(123), "abc")
		p2 := NewPlayer(db, PlayerID(456), "def")
		p3 := NewPlayer(db, PlayerID(789), "ghi")

		r := NewGame(db, p1.ID, p2.ID, p3.ID)
		AssertInt(t, int64(len(r)), 3)

		r, err := p1.Do(db, "11", 0)
		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 3)
	})
}

func TestPlayerCmd(t *testing.T) {
	t.Run("newgame", func(t *testing.T) {
		db := memory.NewStubDatabase()
		
		p1 := NewPlayer(db, PlayerID(123), "abc")

		var err error
		cmd := "/newgame"
		r, err := Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}, p1, 0)
		AssertNoError(t, err)
		_, err = p1.CurrentGame(db)
		AssertNoError(t, err)
		AssertInt(t, int64(len(r)), 1)

	})
	t.Run("leaderboard", func(t *testing.T) {
		db := memory.NewStubDatabase()
		// 
		p1 := NewPlayer(db, PlayerID(123), "abc")

		var err error
		cmd := "/leaderboard"
		r, err := Cmd(db, &tgbotapi.Message{Text: cmd, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd)},
		}}, p1, 0)
		AssertExactError(t, err, NoConnectionError{})
		AssertString(t, r[0].Text, NoConnectionError{}.Error())

		_, err = p1.CurrentGame(db)
		AssertExactError(t, err, NoCurrentGameError{})

	})
	t.Run("no such command", func(t *testing.T) {
		db := memory.NewStubDatabase()
		
		p1 := NewPlayer(db, PlayerID(123), "abc")

		cmd1 := "/command"
		r, err := Cmd(db, &tgbotapi.Message{Text: cmd1, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd1)},
		}}, p1, 0)
		AssertExactError(t, err, NoSuchCommandError{"command"})
		AssertString(t, r[0].Text, NoSuchCommandError{"command"}.Error())

		cmd2 := "/newgame2"
		r, err = Cmd(db, &tgbotapi.Message{Text: cmd2, Entities: []tgbotapi.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: len(cmd2)},
		}}, p1, 0)
		AssertExactError(t, err, NoSuchCommandError{"newgame2"})
		AssertInt(t, int64(len(r)), 1)
		AssertString(t, r[0].Text, NoSuchCommandError{"newgame2"}.Error())

		_, err = p1.CurrentGame(db)
		AssertExactError(t, err, NoCurrentGameError{})
	})

	t.Run("player do new game", func(t *testing.T) {
		db := memory.NewStubDatabase()
		p1 := NewPlayer(db, PlayerID(123), "abc")

		r, err := doNewGame(db, p1, "/newgame")
		AssertNoError(t, err)
		AssertString(t, r[0].Text, "@123 \nStarted\n")
	})
	t.Run("cmd to players id", func(t *testing.T) {
		db := memory.NewStubDatabase()
		_ = NewPlayer(db, PlayerID(123), "abc")
		_ = NewPlayer(db, PlayerID(456), "def")

		players, err := cmdToPlayersID(db, "/newgame @abc @def")
		want := []PlayerID{123, 456}
		AssertNoError(t, err)
		if !reflect.DeepEqual(players, want) {
			t.Errorf("got %v, wanted %v", players, want)
		}
	})
	t.Run("cmd to no players id", func(t *testing.T) {
		db := memory.NewStubDatabase()
		_ = NewPlayer(db, PlayerID(123), "abc")
		_ = NewPlayer(db, PlayerID(456), "def")

		_, err := cmdToPlayersID(db, "/newgame @aaa @bbb")
		AssertExactError(t, err, NoSuchPlayerError{})
	})
}
