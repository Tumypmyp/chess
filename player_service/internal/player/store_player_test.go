package player

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/player_service/pkg/memory"
)

func TestStorePlayer(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := PlayerID(12)
		_, err := getPlayer(id, db)
		AssertExactError(t, err, NoSuchPlayerError{ID: id})

		got := NewPlayer(db, PlayerID(12), "abc")
		want := Player{
			ID:       12,
			Username: "abc",
		}
		AssertPlayer(t, got, want)
	})
	t.Run("store/get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := PlayerID(1234)
		p := Player{ID: id}

		_, err := getPlayer(id, db)
		AssertExactError(t, err, NoSuchPlayerError{ID: id})

		err = StorePlayer(p, db)
		AssertNoError(t, err)

		got, err := getPlayer(id, db)
		AssertNoError(t, err)
		AssertPlayer(t, got, p)
	})
	t.Run("store ID/get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := PlayerID(1234)
		p := Player{
			ID:       id,
			Username: "aba",
		}
		_, err := getID("aba", db)
		AssertExactError(t, err, NoUsernameInDatabaseError{})

		err = StoreID(p, db)
		AssertNoError(t, err)

		got, err := getID("aba", db)
		AssertNoError(t, err)
		AssertPlayerID(t, got, id)
	})
}