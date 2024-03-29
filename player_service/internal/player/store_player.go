package player

import (
	"fmt"

	"github.com/tumypmyp/chess/player_service/pkg/memory"
)

// Update memory.Memory with new value of int64->player
func storePlayer(p Player, m memory.Memory) error {
	key := fmt.Sprintf("user:%d", p.ID)
	if err := m.Set(key, p); err != nil {
		return fmt.Errorf("error when storing player %v: %w", p, err)
	}
	return nil
}

// update memory with new username->int64
func storeID(p Player, m memory.Memory) error {
	key := fmt.Sprintf("username:%s", p.Username)
	if err := m.Set(key, p.ID); err != nil {
		return fmt.Errorf("error when storing player username %v: %w", p.Username, err)
	}
	return nil
}

type NoSuchPlayerError struct {
	ID int64
}

func (n NoSuchPlayerError) Error() string { return fmt.Sprintf("can not get player with id: %v", n.ID) }

// get player from memory
func getPlayer(ID int64, m memory.Memory) (p Player, err error) {
	key := fmt.Sprintf("user:%d", ID)
	if err = m.Get(key, &p); err != nil {
		return p, NoSuchPlayerError{ID: ID}
	}
	return
}

type NoUsernameInDatabaseError struct{}

func (n NoUsernameInDatabaseError) Error() string { return "can not find player by username" }

func getID(username string, db memory.Memory) (id int64, err error) {
	var clientID int64
	key := fmt.Sprintf("username:%v", username)
	if err = db.Get(key, &clientID); err != nil {
		return id, NoUsernameInDatabaseError{}
	}
	return int64(clientID), err
}
