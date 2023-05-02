package player

import (
	"testing"

	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/player_service/pkg/memory"
)

func TestStorePlayer(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(12)
		_, err := getPlayer(id, db)
		AssertExactError(t, err, NoSuchPlayerError{ID: id})

		got := MakePlayer(db, int64(12), "abc")
		want := Player{
			ID:       12,
			Username: "abc",
		}
		AssertPlayer(t, got, want)
	})
	t.Run("store/get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(1234)
		p := Player{ID: id}

		_, err := getPlayer(id, db)
		AssertExactError(t, err, NoSuchPlayerError{ID: id})

		err = storePlayer(p, db)
		AssertNoError(t, err)

		got, err := getPlayer(id, db)
		AssertNoError(t, err)
		AssertPlayer(t, got, p)
	})
	t.Run("store ID/get", func(t *testing.T) {
		db := memory.NewStubDatabase()
		id := int64(1234)
		p := Player{
			ID:       id,
			Username: "aba",
		}
		_, err := getID("aba", db)
		AssertExactError(t, err, NoUsernameInDatabaseError{})

		err = storeID(p, db)
		AssertNoError(t, err)

		got, err := getID("aba", db)
		AssertNoError(t, err)
		AssertInt(t, got, id)
	})
}
