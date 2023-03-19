package memory

import (
	"testing"
	"github.com/tumypmyp/chess/helpers"
)
func TestMemory(t *testing.T) {

// 	t.Run("get", func(t *testing.T) {
// 		db := NewStubDatabase()
// 		var player Player
// 		err := player.Get(PlayerID{12,12}, db)
// 		AssertError(t, err)
// 		player = NewPlayer(db, PlayerID{12,12}, "abc")

// 		if player.ID.ChatID != 12 {
// 			t.Errorf("got %v, want %v", player.ID, PlayerID{12,12})
// 		}
// 	})
// 	t.Run("set/get", func(t *testing.T) {
// 		mem := NewStubDatabase()
// 		p := Player{ID: PlayerID{1234,1234}}
// 		p.Store(mem)
// 		var got Player
// 		got.Get(PlayerID{1234,1234}, mem)
// 		if !reflect.DeepEqual(p, got) {
// 			t.Errorf("got %v, want %v", got, p)
// 		}
// 	})
// }
// func TestStubDatabase(t *testing.T) {
// 	t.Run("string", func(t *testing.T) {
// 		memory := NewStubDatabase()
// 		key := "abcd"
// 		value := "val"
// 		memory.Set(key, value)

// 		var got string
// 		memory.Get(key, &got)
// 		AssertString(t, got, value)
// 	})
// 	t.Run("player", func(t *testing.T) {
// 		memory := NewStubDatabase()
// 		key := "abcd"
// 		value := Player{ID: PlayerID{1234,1234}}
// 		memory.Set(key, value)
// 		t.Logf("mem: %v", *memory.DB)

// 		var got Player
// 		memory.Get(key, &got)
// 		if !reflect.DeepEqual(value, got) {
// 			t.Errorf("got %v, want %v", got, value)
// 		}
// 	})
// 	t.Run("player with game", func(t *testing.T) {
// 		memory := NewStubDatabase()
// 		key := "abcd"
// 		value := Player{ID: PlayerID{1234,1234}, GamesID: []int64{12}}

// 		memory.Set(key, value)

// 		var got Player
// 		memory.Get(key, &got)
// 		if !reflect.DeepEqual(value, got) {
// 			t.Errorf("got %v, want %v", got, value)
// 		}
// 	})


	t.Run("numbers from zero", func(t *testing.T) {
		db := NewStubDatabase()
		val, err := db.Incr("a")
		helpers.AssertNoError(t, err)
		helpers.AssertInt(t, val, 0)
		val, err = db.Incr("a")
		helpers.AssertNoError(t, err)
		helpers.AssertInt(t, val, 1)
		val, err = db.Incr("a")
		helpers.AssertNoError(t, err)
		helpers.AssertInt(t, val, 2)
	})

}
