package main

import (
	"errors"
	"strconv"
)

type Mark int

const (
	Undefined Mark = iota
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
	Player *Player
	Board  [3][3]Mark `json:"board"`
	ID     string
}

func NewGame(p *Player, ID int64) *Game {
	return &Game{
		Player: p,
		ID:     strconv.FormatInt(ID, 10),
	}
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

func (g *Game) move(x, y int) error {
	if x < 0 || len(g.Board) <= x {
		return errors.New("x coordinate out of bounds")
	}
	if y < 0 || len(g.Board[x]) <= y {
		return errors.New("y coordinate out of bounds")
	}
	g.Board[x][y] = First
	return nil
}

func (g *Game) Move(move string) error {
	if len(move) != 2 {
		return errors.New("need 2 characters (example: 22)")
	}
	x := int(move[0] - '0')
	y := int(move[1] - '0')
	if err := g.move(x, y); err != nil {
		return errors.New("bad move: " + err.Error())
	}
	return nil

}
