package main

import (
	"errors"
	"log"
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
	Board [3][3]Mark `json:"board"`
}

func (g *Game) Render() (s string) {
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

func (g *Game) Move(move string) {
	if len(move) != 2 {
		return
	}
	x := int(move[0] - '0')
	y := int(move[1] - '0')
	log.Print(x, y)
	if err := g.move(x, y); err != nil {
		log.Printf("bad move: %v", err)
	}

}
