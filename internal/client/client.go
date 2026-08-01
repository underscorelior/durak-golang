package client

import (
	"durak/internal/game"
)

type ClientList map[*Client]bool

type Client struct {
	ID   uint8
	Name string

	// Provide TableState/Trump/DeckSize without exposing Deck and other Players
}

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
