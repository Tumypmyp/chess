package game

import (
	"testing"
	"github.com/tumypmyp/chess/memory"
	"github.com/tumypmyp/chess/helpers"
)

func TestGame(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := PlayerID(12)
		game := NewGame(db, nil, player)
		err := game.Move(player, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{player},
			Description:   "@12 ",
			CurrentPlayer: 0,
			Status:        Started,
			ChatsID: []int64{12},
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)
	})

	t.Run("2 players", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		p1 := PlayerID(12)
		p2 := PlayerID(13)
		game := NewGame(db, nil, p1, p2)

		err := game.Move(p1, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{p1, p2},
			Description:   "@12 @13 ",
			CurrentPlayer: 1,
			Board:         board,
			
			ChatsID: []int64{12, 13},
			ID:            0}
		AssertGame(t, game, want)

		err = game.Move(p2, "01")
		AssertNoError(t, err)

		board = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		want.Board = board
		want.CurrentPlayer = 0
		AssertGame(t, game, want)
		got := game.String()
		str := "@12 @13 \nStarted\n"
		AssertString(t, got, str)

	})

	t.Run("game status", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := PlayerID(12)
		game := NewGame(db, nil, player)
		AssertStatus(t, game.Status, Started)

		err := game.Move(player, "00")
		AssertNoError(t, err)

		err = game.Move(player, "11")
		AssertNoError(t, err)
		err = game.Move(player, "22")
		AssertNoError(t, err)
		AssertStatus(t, game.Status, Finished)
	})

	t.Run("add chat", func(t *testing.T) {
		db := memory.NewStubDatabase()
		bot := helpers.NewStubBot()
		player := PlayerID(12)
		game := NewGame(db, bot, player)
		AssertStatus(t, game.Status, Started)
		game.AddChat(1234)
		AssertInt(t, int64(len(game.ChatsID)), 2)
		t.Logf("%+v", game.ChatsID)
		game.SendStatus(db, bot)
		AssertInt(t, bot.Len(), 2)
	})
	
}
