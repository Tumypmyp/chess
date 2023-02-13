package main

import (
	"errors"
	"fmt"
)

type Mark int

const (
	UndefinedMark = iota
	First
	Second
)

func (m Mark) String() string {
	switch m {
	case First:
		return "X"
	case Second:
		return "O"
	}
	return "-"
}

type Game struct {
	PlayersID []int64    `json:"players"`
	Board     [3][3]Mark `json:"board"`
	ID        string     `json:"ID:`
}

func NewGame(db Memory, ID string, bot Sender, players ...int64) *Game {
	game := &Game{
		PlayersID: players,
		ID:        ID,
	}
	for _, id := range players {
		var player Player
		db.GetPlayer(id, &player)
		player.SetNewGame(ID)
		db.SetPlayer(player.ID, player)
	}
	db.Set(ID, game)
	game.SendStatus(db, bot)
	return game
}

func (g *Game) String() (s string) {
	for _, row := range g.Board {
		for _, val := range row {
			s += val.String()
		}
		s += "\n"
	}
	return
}

func (g *Game) legalMove(x, y int) (bool, error) {
	if x < 0 || len(g.Board) <= x {
		return false, errors.New("x coordinate out of bounds")
	}
	if y < 0 || len(g.Board[x]) <= y {
		return false, errors.New("y coordinate out of bounds")
	}
	if g.Board[x][y] != UndefinedMark {
		return false, errors.New("this place is not empty")
	}
	return true, nil
}

func (g *Game) findPlayer(id int64) (int, error) {
	for i, p := range g.PlayersID {
		if p == id {
			return i, nil
		}
	}
	return -1, errors.New("player doesnt play this game")
}

func (g *Game) Move(playerID int64, move string) error {
	id, err := g.findPlayer(playerID)
	if err != nil {
		return err
	}

	if len(move) != 2 {
		return errors.New("need 2 characters (example: 22)")
	}
	x := int(move[0] - '0')
	y := int(move[1] - '0')
	if _, err := g.legalMove(x, y); err != nil {
		return fmt.Errorf("illegal move: %w", err)
	}
	g.Board[x][y] = Mark(id + 1)
	return nil
}
func (g *Game) SendStatus(db Memory, bot Sender) {
	for _, id := range g.PlayersID {
		var player Player
		db.GetPlayer(id, &player)
		player.Send(g.String(), bot)
	}
}
