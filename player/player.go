package player

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/leaderboard"
	"github.com/tumypmyp/chess/memory"
	. "github.com/tumypmyp/chess/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (p *Player) NewGame(db memory.Memory, bot Sender, players ...PlayerID) (current_game game.Game) {
	players = append([]PlayerID{p.ID}, players...)

	current_game = game.NewGame(db, players...)

	for _, id := range players {
		p, err := GetPlayer(id, db)
		if err != nil {
			log.Println("no such player", id)
		}
		p.AddNewGame(current_game.ID)
		p.Store(db)
	}

	SendStatus(db, bot, current_game)
	*p, _ = GetPlayer(p.ID, db)
	return
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
	SendStatus(db, bot, game)
	return nil
}

func (p Player) Send(text string, bot Sender) {
	msg := tgbotapi.NewMessage(int64(p.ID), text)

	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}

// sends status to all players
func SendStatus(db memory.Memory, bot Sender, g game.Game) {
	for _, id := range g.ChatsID {
		Send(id, g.String(), makeKeyboard(g), bot)
	}
}
func Send(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup, bot Sender) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	if bot == nil {
		return
	}
	if _, err := bot.Send(msg); err != nil {
		log.Printf("cant send: %v", err)
	}
}

// make inline keyboard for game
func makeKeyboard(g game.Game) tgbotapi.InlineKeyboardMarkup {
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


func  doNewGame(db memory.Memory, bot Sender, p *Player, cmd string) (r Response, err error) {
	others := make([]string, 3)
	n, _ := fmt.Sscanf(cmd, "/newgame @%v @%v @%v", &others[0], &others[1], &others[2])
	others = others[:n]
	var players []PlayerID
	for _, p2 := range others {
		var clientID int64
		key := fmt.Sprintf("username:%v", p2)
		if err = db.Get(key, &clientID); err != nil {
			return Response{}, fmt.Errorf("cant find player @%v", p2)
		}

		id := PlayerID(clientID)

		// var player Player
		// if err := player.Get(id, db); err != nil {
		// 	id.ChatID = clientID
		// 	player.Get(id, db)
		// }
		players = append(players, id)
	}
	g := p.NewGame(db, bot, players...)
	return Response{Text:g.String()}, nil
}

type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

func (p *Player) getLeaderboard(bot Sender) (Response, error) {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: "Vova"})
	if err != nil {
		return Response{Text: NoConnectionError{}.Error()}, NoConnectionError{}

	}
	log.Printf("Greeting: %s", r.GetS())
	// p.Send(r.GetS(), bot)
	return Response{Text: r.GetS()}, nil
}

type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

// runs a command by player
func (p *Player) Cmd(db memory.Memory, bot Sender, cmd *tgbotapi.Message) (r Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	switch cmd.Command() {
	case newgame:
		r, err = doNewGame(db, bot, p, cmd.Text)
	case leaderboard:
		r, err = p.getLeaderboard(bot)
	default:
		err = NoSuchCommandError{cmd.Command()}
		r.Text = err.Error()
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
		return doNewGame(db, bot, p, cmd)
	} else if strings.HasPrefix(cmd, leaderboard) {
		return p.getLeaderboard(bot)
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
