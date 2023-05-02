package player

import (
	"context"
	"fmt"
	"time"
	
	"github.com/tumypmyp/chess/proto/leaderboard"
	pb "github.com/tumypmyp/chess/proto/player"
	"github.com/tumypmyp/chess/player_service/pkg/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NoSuchCommandError struct {
	cmd string
}

func (n NoSuchCommandError) Error() string { return fmt.Sprintf("no such command: %v", n.cmd) }

func NewMessage(p int64, chatID int64, cmd, text string, db memory.Memory) (r pb.Response, err error) {
	if cmd != "" {
		return Cmd(db, cmd, text, p, chatID)
	}
	return Do(p, db, text, chatID)
}

// runs a command by player
func Cmd(db memory.Memory, cmd, text string, p int64, ChatID int64) (r pb.Response, err error) {
	newgame := "newgame"
	leaderboard := "leaderboard"

	// log.Println(cmd, text)
	switch cmd {
	case newgame:
		err = doNewGame(db, p, text)
		r = pb.Response{Text: err.Error(), ChatsID: []int64{ChatID}}
	case leaderboard:
		r, err = getLeaderboard(p, ChatID)
	default:
		err = NoSuchCommandError{cmd}
		r = pb.Response{Text: err.Error(), ChatsID: []int64{ChatID}}
	}
	return
}

// Move player
func Do(id int64, db memory.Memory, move string, chatID int64) (pb.Response, error) {
	// p, _ := getPlayer(id, db)
	game, err := CurrentGame(id, db)
	if err != nil {
		return pb.Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	if err = Move(game, id, move); err != nil {
		return pb.Response{Text: err.Error(), ChatsID: []int64{chatID}}, err
	}
	status, _ := makeStatus(game)
	return pb.Response{Text: status.Description, ChatsID: []int64{int64(id)}}, nil

}


// get or create new player
func MakePlayer(db memory.Memory, playerID int64, username string) (player Player) {
	player, err := getPlayer(playerID, db)
	if err != nil {
		player = NewPlayer(playerID, username)
		storePlayer(player, db)
		storeID(player, db)
	}
	return
}




type NoConnectionError struct{}

func (n NoConnectionError) Error() string { return "can not connect to leaderboard" }

// get leaderboard with gRPC call
func getLeaderboard(id int64, ChatID int64) (pb.Response, error) {
	conn, err := grpc.Dial("leaderboard:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return pb.Response{Text: NoConnectionError{}.Error()}, NoConnectionError{}
	}
	defer conn.Close()
	c := leaderboard.NewLeaderboardClient(conn)

	// Contact the server and return its pb.Response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLeaderboard(ctx, &leaderboard.Player{Name: fmt.Sprintf("%d", id)})
	if err != nil {
		return pb.Response{Text: NoConnectionError{}.Error(), ChatsID: []int64{ChatID}}, NoConnectionError{}

	}
	return pb.Response{Text: r.GetS(), ChatsID: []int64{ChatID}}, nil
}



func doNewGame(db memory.Memory, id int64, cmd string) (err error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return
	}
	players, err := cmdToPlayersID(db, cmd)
	players = append([]int64{p.ID}, players...)
	gameID, err := makeNewGame(players...)
	if err != nil {
		return err
	}
	AddGameToPlayers(db, gameID, players...)

	// _, err = makeStatus(gameID)
	return
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





// returns current game
func CurrentGame(id int64, db memory.Memory) (gameID int64, err error) {
	p, err := getPlayer(id, db)
	if err != nil {
		return
	}
	return p.CurrentGame()
}

// add new game to a player by id
func AddGameToPlayer(db memory.Memory, gameID int64, playerID int64) error {
	p, err := getPlayer(playerID, db)
	if err != nil {
		return err
	}
	p.GamesID = append(p.GamesID, gameID)
	storePlayer(p, db)
	return nil
}

// Add new game
func AddGameToPlayers(db memory.Memory, gameID int64, playersID ...int64) (err error) {
	for _, id := range playersID {
		AddGameToPlayer(db, gameID, id)
	}
	return
}


