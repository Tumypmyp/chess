package player

import (
	"fmt"
	"log"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
	pb "github.com/tumypmyp/chess/proto/player"
)

type Player struct {
	ID       int64
	GamesID  []int64 `json:"gamesID"`
	Username string  `json:"username"`
	Rating   int64
}

// make new player and store in database
func NewPlayer(db memory.Memory, ID int64, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	StorePlayer(p, db)
	StoreID(p, db)
	return p
}

// add new game to a player by id
func AddGameToPlayer(id int64, gameID int64, db memory.Memory) {
	p, err := getPlayer(id, db)
	if err != nil {
		log.Println("no such player", id)
	}
	p.GamesID = append(p.GamesID, gameID)
	StorePlayer(p, db)
}

func (p *Player) addNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}



// make new game
func NewGame(db memory.Memory, players ...int64) pb.Response {
	gameID := makeNewGame(players...)
	for _, id := range players {
		AddGameToPlayer(id, gameID, db)
	}
	status := makeStatus(gameID)
	return pb.Response{Text: status.Description, ChatsID: players}
}



type NoCurrentGameError struct{}

func (n NoCurrentGameError) Error() string { return "no current game,\ntry: /newgame" }

// returns current game
func CurrentGame(id int64, db memory.Memory) (gameID int64, err error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return
	}
	if len(p.GamesID) == 0 {
		return gameID, NoCurrentGameError{}
	}
	gameID = p.GamesID[len(p.GamesID)-1]
	return gameID, nil
}


func cmdToPlayersID(db memory.Memory, cmd string) (playersID []int64, err error) {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]

	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err = db.Get(key, &clientID); err != nil {
			return playersID, NoUsernameInDatabaseError{}
		}

		id := int64(clientID)
		playersID = append(playersID, id)
	}

	return playersID, nil
}

func doNewGame(db memory.Memory, id int64, cmd string) (pb.Response, error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return pb.Response{Text: err.Error()}, err
	}
	players, err := cmdToPlayersID(db, cmd)
	players = append([]int64{p.ID}, players...)
	return NewGame(db, players...), err
}



