package game

import (
	"reflect"
	"testing"

	. "github.com/tumypmyp/chess/helpers"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
)

func TestGameMemory(t *testing.T) {
	t.Run("get/set", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(1234)
		_, err := getGame(id, db)
		AssertExactError(t, err, NoSuchGameError{ID: id})

		game := Game{ID: id}
		err = setGame(game, db)
		AssertNoError(t, err)

		got, err := getGame(id, db)
		AssertNoError(t, err)
		AssertGame(t, got, game)
	})
}

func TestGameMove(t *testing.T) {
	t.Run("new game", func(t *testing.T) {
		db := memory.NewStubDatabase()
		playerID := int64(12)
		gameID := NewGame(db, playerID)
		err := Move(db, gameID, playerID, "00")
		AssertNoError(t, err)

		board := [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}
		want := Game{PlayersID: []int64{playerID},
			Description:   "@ ",
			CurrentPlayer: 0,
			Status:        Started,
			ChatsID:       []int64{12},
			Board:         board,
			ID:            0}
		got, err := getGame(gameID, db)
		AssertNoError(t, err)
		AssertGame(t, got, want)
	})

	t.Run("make status", func(t *testing.T) {
		db := memory.NewStubDatabase()
		playerID := int64(12)
		gameID := NewGame(db, playerID)
		err := Move(db, gameID, playerID, "00")
		AssertNoError(t, err)

		err = Move(db, gameID, playerID, "11")
		AssertNoError(t, err)

		err = Move(db, gameID, playerID, "22")
		AssertNoError(t, err)

		status := MakeStatus(db, gameID)

		got := status.Description
		want := "@ \nFinished\n"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want %+v, got %+v", got, want)
		}
	})
}
