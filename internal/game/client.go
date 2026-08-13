package game

type GamePlayerState struct {
	UserID   string `json:"userId"`
	HandSize int    `json:"handSize"`
}

type GameStateSnapshot struct {
	Players  []GamePlayerState `json:"players"`
	Hand     []Card            `json:"hand"`
	Trump    Card              `json:"trump"`
	DeckSize int               `json:"deckSize"`
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
