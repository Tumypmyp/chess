package main

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tumypmyp/chess/memory"
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
	ID            int64      `json:"ID"`
	Description   string     `json:"description"`
	PlayersID     []PlayerID `json:"players"`
	CurrentPlayer int
	Status        GameStatus
	Board         [3][3]Mark `json:"board"`
}

func NewGame(db memory.Memory, bot Sender, players ...Player) Game {
	ID, err := db.Incr("gameID")
	if err != nil {
		// log.Printf("cant restore id %v", err)
	}
	game := Game{
		ID: ID,
	}
	for _, p := range players {
		game.PlayersID = append(game.PlayersID, p.ID)
		err := p.Get(p.ID, db)
		if err != nil {
			log.Println("no such player", p.ID)
		}
		//log.Println("player", player)
		p.AddNewGame(ID)
		p.Store(db)

		game.Description += "@" + p.Username + " "
	}
	db.Set(fmt.Sprintf("game:%d", ID), game)

	return game
}

// sends status to all players
func (g Game) SendStatus(db memory.Memory, bot Sender) {
	for _, id := range g.PlayersID {
		Send(id, g.String(), bot)
	}
}

func Send(id PlayerID, text string, bot Sender) {
	msg := tgbotapi.NewMessage(id.ChatID, text)

	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}

// Returns string representation of a game
func (g Game) String() (s string) {
	s = g.Description + "\n" + g.Status.String() + "\n"
	for _, row := range g.Board {
		for _, val := range row {
			s += val.String()
		}
		s += "\n"
	}
	return
}

var (
	placeNotEmpty = errors.New("place is not empty")
)

// Returns false if a point out of boundary
func checkBoundary(g Game, x, y int) error {
	if x < 0 || len(g.Board) <= x {
		return errors.New("x coordinate out of bounds")
	}
	if y < 0 || len(g.Board[x]) <= y {
		return errors.New("y coordinate out of bounds")
	}
	if g.Board[x][y] != Undefined {
		return placeNotEmpty
	}
	return nil
}

// Makes move by a player
func (g *Game) Move(playerID PlayerID, move string) error {

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
	if err := checkBoundary(*g, x, y); err != nil {
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
