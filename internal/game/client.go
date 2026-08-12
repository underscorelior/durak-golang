package game

// This describes the visible state of an opponent player to a client
type ClientPlayer struct {
	Name     string `json:"name"`
	UserID   string `json:"userId"`
	Position uint8  `json:"position"`
	HandSize uint8  `json:"handSize"`
}

type ClientGameState struct {
	Players  []ClientPlayer `json:"players"`
	Hand     []Card         `json:"hand"`
	Trump    Card           `json:"trump"`
	DeckSize uint8          `json:"deckSize"`
	Turn     Turn           `json:"turn"`
}
