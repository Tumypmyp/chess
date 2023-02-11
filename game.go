package main

import (
	"errors"
	"fmt"
	"strconv"
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
	Players []*Player  `json:"player"`
	Board   [3][3]Mark `json:"board"`
	ID      string     `json:"ID:`
}

func NewGame(ID int64, players ...*Player) *Game {
	game := &Game{
		Players: players,
		ID:      strconv.FormatInt(ID, 10),
	}
	for _, p := range players {
		p.SetNewGame(game)
	}
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

func (g *Game) findPlayer(player *Player) (int, error) {
	for i, p := range g.Players {
		if p.ChatID == player.ChatID {
			return i, nil
		}
	}
	return -1, errors.New("player doesnt play this game")
}

func (g *Game) Move(player *Player, move string) error {
	id, err := g.findPlayer(player)
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
func (g *Game) SendStatus() {
	for _, p := range g.Players {
		p.Send(g.String())
	}
}
