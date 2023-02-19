package main

import (
	"errors"
	"fmt"
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

type GameStatus int

const (
	Started GameStatus = iota
	Finished
)

func (g GameStatus) String() string {
	switch g {
	case Started:
		return "Started"
	case Finished:
		return "Finished"
	}
	return "unknown"
}

type Game struct {
	ID            int64   `json:"ID"`
	Description   string  `json:"description"`
	PlayersID     []int64 `json:"players"`
	CurrentPlayer int
	Status        GameStatus
	Board         [3][3]Mark `json:"board"`
}

func NewGame(db Memory, bot Sender, players ...int64) Game {
	ID, err := db.Incr("gameID")
	if err != nil {
		// log.Printf("cant restore id %v", err)
	}
	game := Game{
		PlayersID: players,
		ID:        ID,
	}
	for _, id := range players {
		var player Player
		err := player.Get(id, db)
		if err != nil {
			log.Println("no such player", id)
		}
		//log.Println("player", player)
		player.SetNewGame(ID)
		player.Store(db)

		game.Description += "@" + player.Username + " "
	}
	db.Set(fmt.Sprintf("game:%d", ID), game)
	game.SendStatus(db, bot)
	return game
}

func (g *Game) String() (s string) {
	s = g.Description + "\n" + g.Status.String() + "\n"
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
	if g.Board[x][y] != Undefined {
		return false, errors.New("this place is not empty")
	}
	return true, nil
}

func (g *Game) Move(playerID int64, move string) error {

	if g.Status == Finished {
		return errors.New("the game is finished")
	}
	if playerID != g.PlayersID[g.CurrentPlayer] {
		return errors.New("not your turn")
	}

	if len(move) != 2 {
		return errors.New("need 2 characters (example: 22)")
	}
	x := int(move[0] - '0')
	y := int(move[1] - '0')
	if _, err := g.legalMove(x, y); err != nil {
		return fmt.Errorf("illegal move: %w", err)
	}
	g.Board[x][y] = Mark(g.CurrentPlayer + 1)
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.PlayersID)
	g.UpdateStatus()
	return nil
}
func allSame(v [3]Mark) bool {
	return v[0] == v[1] && v[1] == v[2] && v[0] != Undefined
}
func (g *Game) UpdateStatus() {
	if allSame([3]Mark{g.Board[0][0], g.Board[1][1], g.Board[2][2]}) {
		g.Status = Finished
	}
	if allSame([3]Mark{g.Board[2][0], g.Board[1][1], g.Board[0][2]}) {
		g.Status = Finished
	}

	for i := 0; i < 3; i++ {
		if allSame(g.Board[i]) {
			g.Status = Finished
		}
		if allSame([3]Mark{g.Board[2][i], g.Board[1][i], g.Board[0][i]}) {
			g.Status = Finished
		}
	}

}
func (g *Game) SendStatus(db Memory, bot Sender) {
	for _, id := range g.PlayersID {
		var player Player
		player.Get(id, db)
		player.Send(g.String(), bot)
	}
}
