package player

import (
	"log"

	"github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/memory"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)



func Do(update tgbotapi.Update, db memory.Memory, bot game.Sender, cmd string) error {
	player := getPlayerByID(update.SentFrom().ID, update.SentFrom().UserName, db)
	log.Println("player:", player)
	log.Println("message", update.Message)
	if update.Message != nil && update.Message.IsCommand() {
		return player.Cmd(db, bot, update.Message)
	}
	return player.Do(db, bot, cmd)
}

func getPlayerByID(id int64, username string, db memory.Memory) (player Player) {
	ID := game.PlayerID(id)
	var err error
	if player, err = GetPlayer(ID, db); err != nil {
		player = NewPlayer(db, ID, username)
	}
	return
}

func GetPlayer(ID game.PlayerID, m memory.Memory) (p Player, err error) {
	key := fmt.Sprintf("user:%d", ID)
	if err = m.Get(key, &p); err != nil {
		return p, fmt.Errorf("can not get player by id: %w", err)
	}
	return 
}