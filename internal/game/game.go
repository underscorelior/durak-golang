package game

import (
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
)

// TODO: Better name
type PlayerOrder struct {
	Next     string
	Previous string
}
type TableOrder map[string]PlayerOrder

type PlayerState struct {
	hand     []Card
	Position int
}

type Game struct {
	deck    []Card
	players map[string]*PlayerState

	Order TableOrder
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

// FIXME: Lmfao wtf is this name
func (g *Game) FindPlayerWithLowestTrumpCardInHand() string {
	lowestID := ""
	lowestRank := Rank(math.MaxInt)

	trump := g.Trump

	for id, player := range g.players {
		for _, card := range player.hand {
			if card.Suit == trump.Suit {
				if card.Rank == 6 || trump.Rank == 6 && card.Rank == 7 { // TODO: Instead of 6, we need to be able to get the lowest possible card for this optimization
					return id
				}
				if lowestRank > card.Rank {
					lowestID = id
					lowestRank = card.Rank
				}
			}
		}
	}

	if lowestID != "" {
		return lowestID
	}

	// TODO: Pick random player
	return ""
}

func (g *Game) PlayerAtPosition(pos int) string {
	for id, player := range g.players {
		if player.Position == pos {
			return id
		}
	}

	return ""
}

func (g *Game) SetupTableOrder() {
	numPlayers := len(g.players)
	tableOrder := make(TableOrder, numPlayers)

	players := g.players

	keys := slices.Collect(maps.Keys(players))

	sort.Slice(keys[:], func(i, j int) bool {
		return players[keys[i]].Position < players[keys[j]].Position
	})

	for _, key := range keys {
		pos := players[key].Position
		tableOrder[key] = PlayerOrder{
			Next:     g.PlayerAtPosition((pos + 1) % numPlayers),
			Previous: g.PlayerAtPosition((pos - 1) % numPlayers),
		}
	}

	g.Order = tableOrder
}

func (g *Game) SetupInitialTurn() {
	initialAttacker := g.FindPlayerWithLowestTrumpCardInHand()
	defender := g.Order[initialAttacker].Next

	turn := Turn{
		TableState:        make([]CardPair, 6),
		DefenderID:        defender,
		InitialAttackerID: initialAttacker,
		AttackerIDs:       g.AdjacentTo(defender),

		Phase: INITIAL,
	}

	g.Turn = turn
}

// Creates and shuffles the deck, deals cards and picks the Trump card
// TODO: Fix the argument type
func InitializeGame(players map[string]struct{ Position int }) *Game {
	var g Game
	g.players = make(map[string]*PlayerState, len(players))
	for k, v := range players {
		g.players[k] = &PlayerState{Position: v.Position}
	}

	g.deck = CreateDeck()
	g.ShuffleDeck()

	g.DealCards()
	g.Trump, g.deck = g.deck[0], g.deck[1:]

	g.SetupTableOrder()

	return &g
}

func (g *Game) DeckSize() int {
	return len(g.deck)
}

func (g *Game) AdjacentTo(userID string) map[string]bool {
	o := g.Order[userID]
	return map[string]bool{o.Next: true, o.Previous: false}
}
