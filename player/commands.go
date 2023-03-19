package player

import (
	"log"

	"github.com/tumypmyp/chess/game"
	"github.com/tumypmyp/chess/memory"
	
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getPlayerByID(id int64, username string, db memory.Memory) (player Player) {
	ID := game.PlayerID(id)
	if err := player.Get(ID, db); err != nil {
		player = NewPlayer(db, ID, username)
	}
	return
}

func Do(update tgbotapi.Update, db memory.Memory, bot game.Sender, cmd string) error {
	player := getPlayerByID(update.SentFrom().ID, update.SentFrom().UserName, db)
	log.Println("player:", player)
	log.Println("message", update.Message)
	if update.Message != nil && update.Message.IsCommand() {
		return player.Cmd(db, bot, update.Message)
	}
	return player.Do(db, bot, cmd)
}
