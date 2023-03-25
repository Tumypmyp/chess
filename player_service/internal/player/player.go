package player

import (
	"fmt"
	"log"

	"github.com/tumypmyp/chess/player_service/internal/game"
	. "github.com/tumypmyp/chess/helpers"
	
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/player"
)

type Player struct {
	ID       PlayerID
	GamesID  []int64 `json:"gamesID"`
	Username string  `json:"username"`
}

func NewPlayer(db memory.Memory, ID PlayerID, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	db.Set(fmt.Sprintf("username:%v", p.Username), p.ID)
	Store(p, db)
	StoreUsername(p, db)
	return p
}

type NoCurrentGameError struct{}

func (n NoCurrentGameError) Error() string { return "no current game,\ntry: /newgame" }

func (p *Player) CurrentGame(db memory.Memory) (game game.Game, err error) {
	*p, err = getPlayer(p.ID, db)
	if err != nil {
		return
	}
	if len(p.GamesID) == 0 {
		return game, NoCurrentGameError{}
	}
	err = db.Get(fmt.Sprintf("game:%d", p.GamesID[len(p.GamesID)-1]), &game)
	return
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}

func NewGame(db memory.Memory, players ...PlayerID) pb.Response {
	current_game := game.NewGame(db, players...)

	for _, id := range players {
		p, err := getPlayer(id, db)
		if err != nil {
			log.Println("no such player", id)
		}
		p.AddNewGame(current_game.ID)
		Store(p, db)
	}

	return SendStatus(current_game)
}

type NoUsernameInDatabaseError struct{}

func (n NoUsernameInDatabaseError) Error() string { return "can not find player by username" }

func getID(username string, db memory.Memory) (id PlayerID, err error) {
	var clientID int64
	key := fmt.Sprintf("username:%v", username)
	if err = db.Get(key, &clientID); err != nil {
		return id, NoUsernameInDatabaseError{}
	}
	return PlayerID(clientID), err
}


func cmdToPlayersID(db memory.Memory, cmd string) (playersID []PlayerID, err error) {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]

	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err = db.Get(key, &clientID); err != nil {
			return playersID, NoUsernameInDatabaseError{}
		}

		id := PlayerID(clientID)
		playersID = append(playersID, id)
	}

	return playersID, nil
}

func doNewGame(db memory.Memory, id PlayerID, cmd string) (pb.Response, error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return pb.Response{Text: err.Error()}, err
	}
	players, err := cmdToPlayersID(db, cmd)
	players = append([]PlayerID{p.ID}, players...)
	return NewGame(db, players...), err
}


// sends status to all players
func SendStatus(g game.Game) pb.Response {
	return pb.Response{Text: g.String(), Keyboard: makeGameKeyboard(g), ChatsID: g.ChatsID}
}

func makeGameKeyboard(g game.Game) (keyboard []*pb.ArrayButton) {
	keyboard = make([]*pb.ArrayButton, len(g.Board))

	for i, v := range g.Board {
		keyboard[i] = &pb.ArrayButton{Buttons: make([]*pb.Button, len(v))}
		for j, b := range v {
			keyboard[i].Buttons[j] = &pb.Button{Text: b.String(), CallbackData: fmt.Sprintf("%d%d", i, j)}
		}
	}
	
	// log.Println(keyboard)
	return
}

// Update memory.Memory with new value of a player
func Store(p Player, m memory.Memory) error {
	key := fmt.Sprintf("user:%d", p.ID)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}

func StoreUsername(p Player, m memory.Memory) error {
	key := fmt.Sprintf("username:%s", p.Username)
	if err := m.Set(key, p.ID); err != nil {
		return fmt.Errorf("error when storing player username %v: %w", p.Username, err)
	}
	return nil
}
