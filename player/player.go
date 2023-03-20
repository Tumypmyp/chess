package player

import (
	"fmt"
	"log"
	"strings"

	"github.com/tumypmyp/chess/game"
	. "github.com/tumypmyp/chess/helpers"
	"github.com/tumypmyp/chess/memory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	*p, err = GetPlayer(p.ID, db)
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

func NewGame(db memory.Memory, players ...PlayerID) []Response {
	current_game := game.NewGame(db, players...)

	for _, id := range players {
		p, err := GetPlayer(id, db)
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

func doNewGame(db memory.Memory, p *Player, cmd string) ([]Response, error) {
	players, err := cmdToPlayersID(db, cmd)
	players = append([]PlayerID{p.ID}, players...)
	return NewGame(db, players...), err
}


// add p.Update()

func (p *Player) Move(db memory.Memory, bot Sender, move string) error {
	game, err := p.CurrentGame(db)
	if err != nil {
		return err
	}
	if err = game.Move(p.ID, move); err != nil {
		return err
	}
	if err := db.Set(fmt.Sprintf("game:%d", game.ID), game); err != nil {
		return fmt.Errorf("could not reach db: %w", err)
	}
	SendStatus(game)
	return nil
}


// sends status to all players
func SendStatus(g game.Game) (r []Response) {
	for _, _ = range g.ChatsID {
		r = append(r, Response{Text:g.String(), Keyboard: makeGameKeyboard(g)})
	}
	return
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



type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

// runs a command by player
func (p *Player) Cmd(db memory.Memory, bot Sender, cmd *tgbotapi.Message) (r []Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	switch cmd.Command() {
	case newgame:
		r, err = doNewGame(db, p, cmd.Text)
	case leaderboard:
		r1, err2 := getLeaderboard(*p)
		r = []Response{r1}
		err = err2
	default:
		err = NoSuchCommandError{cmd.Command()}
		r = []Response{Response{Text:err.Error()}}
	}
	return
}

func (p *Player) Do(db memory.Memory, bot Sender, cmd string) (r Response, err error) {
	r, err = p.Do2(db, bot, cmd)
	if err != nil {
		return Response{Text: err.Error()}, err
		// 	p.Send(err.Error(), bot)
	}
	return
}

func (p *Player) Do2(db memory.Memory, bot Sender, cmd string) (Response, error) {
	pref := "/newgame"
	leaderboard := "/leaderboard"

	if strings.HasPrefix(cmd, pref) {
		return Response{}, nil
	} else if strings.HasPrefix(cmd, leaderboard) {
		return getLeaderboard(*p)
	} else {
		return Response{}, p.Move(db, bot, cmd)
	}
}

// Update memory.Memory with new value of a player
func (p Player) Store(m memory.Memory) error {
	key := fmt.Sprintf("user:%d", p.ID)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}
