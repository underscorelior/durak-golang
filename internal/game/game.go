package game

import (
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
)

type PlayerState struct {
	hand     []Card
	Position int
	Passed   bool
}

type Game struct {
	deck    []Card
	players map[string]*PlayerState

	Order TableOrder
	Trump Card
	Turn  Turn

	FirstDefense bool
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
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
		g.players[i].hand = slices.Clone(g.deck[:6]) // TODO: Fix magic number : 6
		g.deck = g.deck[6:]
	}
} // TODO: Make a test that covers when two pplayers deal and then draw to the first player's hand then ensure that the second player's hand is unchanged (ensuring that the data is a copy and not a reference)

func (g *Game) FindPlayerWithLowestTrumpCardInHand() string {
	lowestID := ""
	lowestRank := Rank(math.MaxInt)

	trump := g.Trump

	for id, player := range g.players {
		for _, card := range player.hand {
			if card.Suit == trump.Suit {
				if card.Rank == 6 || trump.Rank == 6 && card.Rank == 7 {
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

	i := rand.IntN(len(g.players))
	return slices.Collect(maps.Keys(g.players))[i]
}

func (g *Game) StartNextTurn(initial bool, skipDefender bool) {
	var initialAttacker string

	if initial {
		initialAttacker = g.FindPlayerWithLowestTrumpCardInHand()
	} else if skipDefender {
		initialAttacker = g.Order[g.Turn.DefenderID].Next
	} else {
		initialAttacker = g.Turn.DefenderID
	}

	defender := g.Order[initialAttacker].Next

	turn := Turn{
		TableState:        make([]CardPair, 6),
		DefenderID:        defender,
		InitialAttackerID: initialAttacker,
		AttackerIDs:       g.AdjacentTo(defender),

		Phase: TurnPhaseInitial,
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
	g.Trump = g.deck[0]
	g.deck = append(g.deck[1:], g.Trump)

	g.SetupTableOrder()

	g.StartNextTurn(true, false)

	return &g
}

func (g *Game) DeckSize() int {
	return len(g.deck)
}

// Draws specified number of cards into player's hand. Returns a boolean representing if any more cards can be drawn.
func (g *Game) DrawCards(userID string, amount int) bool {
	p := g.players[userID]

	if p == nil {
		// TODO: Some sort of error handling
		return true
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

func (p *PlayerState) RemoveCard(card Card) []Card {
	hand := slices.Collect(func(yield func(Card) bool) {
		for _, c := range p.hand {
			if c.Rank != card.Rank || c.Suit != card.Suit {
				if !yield(c) {
					return
				}
			}
		}
	})

	return hand
}

func (p *PlayerState) CardInHand(card Card) bool {
	return slices.Contains(p.hand, card)
}
