package game

import (
	"math/rand/v2"
	"slices"
)

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
	for _, player := range g.players {
		player.hand = slices.Clone(g.deck[:6]) // TODO: Fix magic number : 6
		g.deck = g.deck[6:]
	}
} // TODO: Make a test that covers when two pplayers deal and then draw to the first player's hand then ensure that the second player's hand is unchanged (ensuring that the data is a copy and not a reference)

func (g *Game) DeckSize() int {
	return len(g.deck)
}

// TODO: Overhaul this function to make it more logical, this and DistributeCards()
// Draws specified number of cards into player's hand. Returns a boolean representing if any more cards can be drawn.
func (g *Game) DrawCards(userID string, amount int) bool {
	p := g.players[userID]

	if p == nil {
		// TODO: Some sort of error handling
		return true
	}

	if len(g.deck) <= 0 {
		return false
	}

	for range max(amount, 0) {
		p.hand = append(p.hand, g.deck[0])
		g.deck = g.deck[1:]

		// TODO: NEED TO ACCOUNT FOR TRUMP CARD
		if len(g.deck) <= 0 {
			return false
		}
	}

	return true
}

func (g *Game) DistributeCards() {
	ia := g.Turn.InitialAttackerID

	g.DrawCards(ia, 6-len(g.players[ia].hand))

	for attacker := range g.Turn.AttackerIDs {
		if attacker != ia {
			g.DrawCards(attacker, 6-len(g.players[attacker].hand))
		}
	}

	g.DrawCards(g.Turn.DefenderID, 6-len(g.players[g.Turn.DefenderID].hand))
}
