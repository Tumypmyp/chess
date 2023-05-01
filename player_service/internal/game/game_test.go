package game

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"
	
	"github.com/tumypmyp/chess/player_service/pkg/memory"
)



func TestGame(t *testing.T) {
	t.Run("makeMove", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := PlayerID(12)
		game := makeGame(db, player)
		err := game.makeMove(player, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{player},
			Description:   "@ ",
			CurrentPlayer: 0,
			Status:        Started,
			ChatsID:       []int64{12},
			Board:         board,
			ID:            0}
		AssertGame(t, game, want)
	})

	t.Run("2 players", func(t *testing.T) {
		mem := memory.NewStubDatabase()
		db := mem
		p1 := PlayerID(12)
		p2 := PlayerID(13)
		game := makeGame(db, p1, p2)

		err := game.makeMove(p1, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []PlayerID{p1, p2},
			Description:   "@ @ ",
			CurrentPlayer: 1,
			Board:         board,

			ChatsID: []int64{12, 13},
			ID:      0}
		AssertGame(t, game, want)

		err = game.makeMove(p2, "01")
		AssertNoError(t, err)

		board = [3][3]Mark{{1, 2, 0}, {0, 0, 0}, {0, 0, 0}}
		want.Board = board
		want.CurrentPlayer = 0
		AssertGame(t, game, want)
		got := game.String()
		str := "@ @ \nStarted\n"
		AssertString(t, got, str)

	})

	t.Run("game status", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := PlayerID(12)
		game := makeGame(db, player)
		AssertStatus(t, game.Status, Started)

		err := game.makeMove(player, "00")
		AssertNoError(t, err)

		err = game.makeMove(player, "11")
		AssertNoError(t, err)
		err = game.makeMove(player, "22")
		AssertNoError(t, err)
		AssertStatus(t, game.Status, Finished)
	})

	t.Run("add chat", func(t *testing.T) {
		db := memory.NewStubDatabase()
		player := PlayerID(12)
		game := makeGame(db, player)
		AssertStatus(t, game.Status, Started)
		game.AddChat(1234)
		AssertInt(t, int64(len(game.ChatsID)), 2)
		t.Logf("%+v", game.ChatsID)
	})

}
