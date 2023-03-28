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
	Rating   int64
}

// make new player and store in database
func NewPlayer(db memory.Memory, ID PlayerID, Username string) Player {
	p := Player{
		ID:       ID,
		Username: Username,
	}
	StorePlayer(p, db)
	StoreID(p, db)
	return p
}

// add new game to a player by id
func AddGameToPlayer(id PlayerID, gameID int64, db memory.Memory) {
	p, err := getPlayer(id, db)
	if err != nil {
		log.Println("no such player", id)
	}
	p.GamesID = append(p.GamesID, gameID)
	StorePlayer(p, db)
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}



// make new game
func NewGame(db memory.Memory, players ...PlayerID) pb.Response {
	g := game.NewGame(db, players...)
	for _, id := range players {
		AddGameToPlayer(id, g.ID, db)
	}
	return game.SendStatus(g)
}



type NoCurrentGameError struct{}

func (n NoCurrentGameError) Error() string { return "no current game,\ntry: /newgame" }


func CurrentGame(id PlayerID, db memory.Memory) (g game.Game, err error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return g, err
	}
	if len(p.GamesID) == 0 {
		return g, NoCurrentGameError{}
	}
	gameID := p.GamesID[len(p.GamesID)-1]
	g, err = game.GetGame(gameID, db)
	return
}

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



