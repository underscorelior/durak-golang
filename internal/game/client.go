package game

// This describes the visible state of an opponent player to a client
type ClientPlayer struct {
	Name string
	ID   uint8

	Position uint8
	HandSize uint8
}

type ClientGameState struct {
	Players []ClientPlayer
	hand    []Card

	Trump    Card
	DeckSize uint8

	Turn Turn
}
