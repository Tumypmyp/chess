package game

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

// Representation of a play
type Game struct {
	// game id
	ID int64 `json:"ID"`
	// string describing a game status, players
	Description string `json:"description"`
	// players of a game
	PlayersID []PlayerID `json:"players"`
	// player id in a slice
	CurrentPlayer int
	// status of a game
	Status GameStatus
	// board representation
	Board [3][3]Mark `json:"board"`
}

func NewGame(db memory.Memory, bot Sender, players ...PlayerID) Game {
	ID, err := db.Incr("gameID")
	if err != nil {
		// log.Printf("cant restore id %v", err)
	}
	game := Game{
		ID: ID,
	}
	// make description
	for _, id := range players {
		game.PlayersID = append(game.PlayersID, id)

		game.Description += "@" + string(id) + " "
	}

	db.Set(fmt.Sprintf("game:%d", ID), game)
	return game
}

// sends status to all players
func (g Game) SendStatus(db memory.Memory, bot Sender) {
	for _, id := range g.PlayersID {
		Send(id, g.String(), makeKeyboard(g), bot)
	}
}

// make inline keyboard for game
func makeKeyboard(g Game) tgbotapi.InlineKeyboardMarkup {
	markup := make([][]tgbotapi.InlineKeyboardButton, len(g.Board))

	for i, v := range g.Board {
		markup[i] = make([]tgbotapi.InlineKeyboardButton, len(g.Board[i]))
		for j, _ := range v {
			markup[i][j] = tgbotapi.NewInlineKeyboardButtonData(g.Board[i][j].String(), fmt.Sprintf("%d%d", i, j))
		}
	}
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: markup,
	}
}

func Send(id PlayerID, text string, keyboard tgbotapi.InlineKeyboardMarkup, bot Sender) {
	msg := tgbotapi.NewMessage(int64(id), text)
	msg.ReplyMarkup = keyboard
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
	return
}

type PlaceNotEmptyError struct{}

func (n PlaceNotEmptyError) Error() string { return "place is not empty" }

// Returns false if a point out of boundary
func checkBoundary(g Game, x, y int) error {
	if x < 0 || len(g.Board) <= x {
		return errors.New("x coordinate out of bounds")
	}
	if y < 0 || len(g.Board[x]) <= y {
		return errors.New("y coordinate out of bounds")
	}
	if g.Board[x][y] != Undefined {
		return PlaceNotEmptyError{}
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

// update status of a game
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
