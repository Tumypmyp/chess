package main

import (
	"errors"
	"fmt"
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
	Player *Player    `json:"player"`
	Board  [3][3]Mark `json:"board"`
	ID     string     `json:"ID:`
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

func (g *Game) legalMove(x, y int) (bool, error) {
	if x < 0 || len(g.Board) <= x {
		return false, errors.New("x coordinate out of bounds")
	}
	if y < 0 || len(g.Board[x]) <= y {
		return false, errors.New("y coordinate out of bounds")
	}
	return true, nil
}

func (g *Game) Move(move string) error {
	if len(move) != 2 {
		return errors.New("need 2 characters (example: 22)")
	}
	x := int(move[0] - '0')
	y := int(move[1] - '0')
	if _, err := g.legalMove(x, y); err != nil {
		return fmt.Errorf("illegal move: %w", err)
	}
	g.Board[x][y] = First
	return nil
}
