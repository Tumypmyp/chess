package game

import (
	"fmt"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/game"
)

func NewGame(db memory.Memory, playersID ...int64) int64 {
	return makeGame(db, playersID...).ID	
}

type NoSuchGameError struct {
	ID int64
}

func (n NoSuchGameError) Error() string { return fmt.Sprintf("can not get game with id: %v", n.ID) }

// get game from memory
func getGame(ID int64, m memory.Memory) (p Game, err error) {
	key := fmt.Sprintf("game:%d", ID)
	if err = m.Get(key, &p); err != nil {
		return p, NoSuchGameError{ID: ID}
	}
	return
}

// Update memory.Memory with new value of gameID->game
func setGame(g Game, m memory.Memory) error {
	key := fmt.Sprintf("game:%d", g.ID)
	if err := m.Set(key, g); err != nil {
		return fmt.Errorf("error when storing game %v: %w", g, err)
	}
	return nil
}

// make move in game for player
func Move(m memory.Memory, gameID int64, playerID int64, move string) (err error) {
	g, err := getGame(gameID, m)
	if err != nil {
		return err
	}
	g.makeMove(playerID, move)
	if err = setGame(g, m); err != nil {
		return err
	}
	return
}



// make status
func MakeStatus(m memory.Memory, gameID int64) pb.GameStatus {
	g, _ := getGame(gameID, m)
	return pb.GameStatus{Description: g.String(), Keyboard: makeGameKeyboard(g)}
}

func makeGameKeyboard(g Game) (keyboard []*pb.ArrayButton) {
	keyboard = make([]*pb.ArrayButton, len(g.Board))

	for i, v := range g.Board {
		keyboard[i] = &pb.ArrayButton{Buttons: make([]*pb.Button, len(v))}
		for j, b := range v {
			keyboard[i].Buttons[j] = &pb.Button{Text: b.String(), CallbackData: fmt.Sprintf("%d%d", i, j)}
		}
	}
	return
}