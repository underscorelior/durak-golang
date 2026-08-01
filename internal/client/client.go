package client

import (
	"durak/internal/game"
)

// This describes the visible state of an opponent player to a client
type ClientPlayer struct {
	Name string
	ID   uint8

	Position uint8
	HandSize uint8
}

type ClientGameState struct {
	Players []ClientPlayer
	hand    []game.Card

	Trump    game.Card
	DeckSize uint8

	Turn game.Turn
}
