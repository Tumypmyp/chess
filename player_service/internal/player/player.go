package player

import (
	"fmt"
	"log"

	"github.com/tumypmyp/chess/player_service/internal/game"
	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"
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
	p.Store(db)
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

func NewGame(db memory.Memory, players ...PlayerID) Response {
	current_game := game.NewGame(db, players...)

	for _, id := range players {
		p, err := getPlayer(id, db)
		if err != nil {
			log.Println("no such player", id)
		}
		p.AddNewGame(current_game.ID)
		p.Store(db)
	}

	return SendStatus(current_game)
}

type NoSuchPlayerError struct{}

func (n NoSuchPlayerError) Error() string { return "can not find player" }

func cmdToPlayersID(db memory.Memory, cmd string) (playersID []PlayerID, err error) {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]

	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err = db.Get(key, &clientID); err != nil {
			return playersID, NoSuchPlayerError{}
		}

		id := PlayerID(clientID)
		playersID = append(playersID, id)
	}

	return playersID, nil
}

func doNewGame(db memory.Memory, id PlayerID, cmd string) (Response, error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return Response{Text: err.Error()}, err
	}
	players, err := cmdToPlayersID(db, cmd)
	players = append([]PlayerID{p.ID}, players...)
	return NewGame(db, players...), err
}

// add p.Update()

// sends status to all players
func SendStatus(g game.Game) Response {
	return Response{Text: g.String(), Keyboard: makeGameKeyboard(g), ChatsID: g.ChatsID}
}

func makeGameKeyboard(g game.Game) (keyboard [][]Button) {
	keyboard = make([][]Button, len(g.Board))

	for i, v := range g.Board {
		keyboard[i] = make([]Button, len(g.Board[i]))
		for j, _ := range v {
			keyboard[i][j] = Button{g.Board[i][j].String(), fmt.Sprintf("%d%d", i, j)}
		}
	}
	return
}

// Update memory.Memory with new value of a player
func (p Player) Store(m memory.Memory) error {
	key := fmt.Sprintf("user:%d", p.ID)
	p.StoreID(m)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}

func (p Player) StoreID(m memory.Memory) error {
	key := fmt.Sprintf("userID:%d", p.ID)
	if err := m.Set(key, p.Username); err != nil {
		return fmt.Errorf("error when storing player username %v: %w", p.Username, err)
	}
	return nil
}
