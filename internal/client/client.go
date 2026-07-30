package client

import (
	"durak/internal/game"
)

type Client struct {
	ID   uint8
	Name string

	// Provide TableState/Trump/DeckSize without exposing Deck and other Players
}

type ClientPlayer struct {
	ID   uint8
	Name string

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
