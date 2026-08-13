package game

import (
	"math/rand/v2"
)

type PlayerState struct {
	hand []Card
}

type Game struct {
	deck    []Card
	players map[string]*PlayerState

	Trump Card
	Turn  Turn
}

// Generates a full deck of cards
func CreateDeck() []Card {
	var deck []Card
	for _, suit := range Suits {
		for _, rank := range Ranks {
			deck = append(deck, Card{suit, rank})
		}
	}

	return deck
}

// Shuffles a deck (in-place)
func (g *Game) ShuffleDeck() {
	rand.Shuffle(len(g.deck), func(i, j int) {
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	})
}

// Deals cards to the Players in the game, removes those cards from the deck
func (g *Game) DealCards() {
	for i := range g.players {
		g.players[i].hand = g.deck[:6] // TODO: Fix magic number : 6
		g.deck = g.deck[6:]
	}
}

// Creates and shuffles the deck, deals cards and picks the Trump card
func (g *Game) InitializeGame() {
	g.deck = CreateDeck()
	g.ShuffleDeck()

	g.DealCards()
	g.Trump, g.deck = g.deck[0], g.deck[1:]
}

func (g *Game) DeckSize() int {
	return len(g.deck)
}
