package player

type Player struct {
	ID       int64
	GamesID  []int64 `json:"gamesID"`
	Username string  `json:"username"`
	Rating   int64
}

func NewPlayer(ID int64, Username string) Player {
	return Player{
		ID:       ID,
		Username: Username,
	}
}

func (p *Player) AddNewGame(gameID int64) {
	p.GamesID = append(p.GamesID, gameID)
}


type NoGamesError struct{}

func (n NoGamesError) Error() string { return "no games started" }

func (p *Player) CurrentGame() (gameID int64, err error) {
	if len(p.GamesID) == 0 {
		return gameID, NoGamesError{}
	}
	gameID = p.GamesID[len(p.GamesID)-1]
	return gameID, nil
}
