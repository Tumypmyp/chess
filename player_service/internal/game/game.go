package game

import (
	"errors"
	"fmt"
	"log"

	. "github.com/tumypmyp/chess/helpers"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
	
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
	// chats to send messages
	ChatsID []int64
	// status of a game
	Status GameStatus
	// board representation
	Board [3][3]Mark `json:"board"`
}

func makeGame(db memory.Memory, playersID ...PlayerID) Game {
	ID, err := db.Incr("gameID")
	if err != nil {
		log.Printf("cant restore id: %v", err)
	}
	game := Game{
		ID: ID,
	}
	// makes description
	for _, id := range playersID {
		game.PlayersID = append(game.PlayersID, id)
		game.ChatsID = append(game.ChatsID, int64(id))
		name, _ := getPlayerUsername(id, db)
		game.Description += fmt.Sprintf("@%s ", name)
	}

	setGame(game, db)
	return game
}


// get player from memory
func getPlayerUsername(ID PlayerID, m memory.Memory) (string, error) {
	key := fmt.Sprintf("user:%d", ID)
	var p struct {
		ID       PlayerID
		GamesID  []int64 `json:"gamesID"`
		Username string  `json:"username"`
	}
	if err := m.Get(key, &p); err != nil {
		return p.Username, fmt.Errorf("can not get player by id: %w", err)
	}
	return p.Username, nil
}

func (g *Game) AddChat(chatID int64) {
	g.ChatsID = append(g.ChatsID, chatID)
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
func (g *Game) makeMove(playerID PlayerID, move string) error {
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
