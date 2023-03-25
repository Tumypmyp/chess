package memory

import (
	"reflect"
	"testing"

	. "github.com/tumypmyp/chess/helpers"
)

func TestMemory(t *testing.T) {

}
func TestStubDatabase(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		memory := NewStubDatabase()
		key := "abcd"
		value := "val"
		memory.Set(key, value)

		var got string
		memory.Get(key, &got)
		AssertString(t, got, value)
	})
	t.Run("struct", func(t *testing.T) {
		memory := NewStubDatabase()
		key := "abcd"
		value := struct {
			F float64
			I int
			S string
			A []bool
			p int64
		}{
			F: 3.14,
			I: 123,
			S: "abcd",
			A: []bool{true, false, false},
			p: 111,
		}

		memory.Set(key, value)

		var got struct {
			F float64
			I int
			S string
			A []bool
			p int64
		}
		want := value
		want.p = 0

		memory.Get(key, &got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

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

//			var got Player
//			memory.Get(key, &got)
//			if !reflect.DeepEqual(value, got) {
//				t.Errorf("got %v, want %v", got, value)
//			}
//		})
//	}
func TestMemoryIncr(t *testing.T) {
	t.Run("numbers from zero", func(t *testing.T) {
		db := NewStubDatabase()
		val, err := db.Incr("a")
		AssertNoError(t, err)
		AssertInt(t, val, 0)
		val, err = db.Incr("a")
		AssertNoError(t, err)
		AssertInt(t, val, 1)
		val, err = db.Incr("a")
		AssertNoError(t, err)
		AssertInt(t, val, 2)
	})

}
