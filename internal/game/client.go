package game

type GamePlayerState struct {
	UserID   string `json:"user_id"`
	HandSize int    `json:"hand_size"`
}

type GameStateSnapshot struct {
	Players  []GamePlayerState `json:"players"`
	Hand     []Card            `json:"hand"`
	Trump    Card              `json:"trump"`
	DeckSize int               `json:"deck_size"`
	Turn     Turn              `json:"turn"`
}

func (g *Game) playerSnapshots() []GamePlayerState {
	var playerSnapshots []GamePlayerState
	for userID, player := range g.players {
		playerSnapshot := GamePlayerState{
			UserID:   userID,
			HandSize: len(player.hand),
		}
		playerSnapshots = append(playerSnapshots, playerSnapshot)
	}

	return playerSnapshots
}

func (g *Game) StateFor(userID string) *GameStateSnapshot {
	return &GameStateSnapshot{
		Players:  g.playerSnapshots(),
		Hand:     g.players[userID].hand,
		Trump:    g.Trump,
		DeckSize: g.DeckSize(),
		Turn:     g.Turn,
	}
}
