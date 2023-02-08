package main

import (
	"reflect"
	"testing"
)

func TestGameMove(t *testing.T) {
	t.Run("move", func(t *testing.T) {
		game := Game{}
		game.Move("00")
		want := Game{Board: [3][3]Mark{{1, 0, 0}, {0, 0, 0}, {0, 0, 0}}}
		if !reflect.DeepEqual(game, want) {
			t.Fatalf("wanted %v, got %v", want, game)
		}
	})
}
