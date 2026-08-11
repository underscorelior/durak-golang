package game

// This describes the visible state of an opponent player to a client
type ClientPlayer struct {
	Name    string
	LobbyID string // I think this should be a "LobbyID" where the other ID is the "GeneralID" (one is the general one, used across the manager, one is lobby only)

	Position uint8 // Need a new pattern for this
	HandSize uint8
}

type ClientGameState struct {
	Players []ClientPlayer
	hand    []Card

	Trump    Card
	DeckSize uint8

	Turn Turn
}
