package game

import (
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
)

type Game struct {
	deck    []Card
	players map[string]*PlayerState

	Order TableOrder
	Trump Card
	Turn  Turn

	FirstDefense bool
}

// TODO: Find a better way to do client/server err handling
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}

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

func (g *Game) StartNextTurn(initialAttacker string) {
	g.ResetPassingState()

	defender := g.Order[initialAttacker].Next

	tableSize := 6
	if !g.FirstDefense {
		tableSize = 5
	}

	turn := Turn{
		TableState:        make([]CardPair, tableSize),
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

	g.StartNextTurn(g.FindPlayerWithLowestTrumpCardInHand())

	return &g
}

func (g *Game) EndTurn(endPhase TurnPhase) {

}

func (g *Game) MakeAction(userID string, action Action, card *Card, pairIndex *int) error {
	t := g.Turn
	p := g.players[userID]

	if p == nil {
		return Error{
			Code:    "INVALID_USER",
			Message: "Provided user ID is not in the game",
		}
	}

	// TAKE: Only DEFENDER during GAME phase
	if action == ActionTake {
		if err := g.CanTake(userID); err != nil {
			return err
		}

		// START TIMER ALLOWING ATTACKERS TO THROW IN
		// NEED A WAY TO ACCURATELY DISPLAY IT (SENDING TIMESTAMP AND ALLOWING CLIENT TO CALC HOW LONG IS LEFT?)
		// ONCE TIME IS UP, MOVE CARDS TO PLAYER HAND
		// p.hand = append(p.hand, t.GetPlayedCards()...)

		return nil
	}

	// PASS: Only ATTACKERs during GAME phase, PASS BY DEFAULT ONCE CANNOT ADD NEW CARDS, PASSING CANCELLED ONCE A NEW CARD IS PLACED
	if action == ActionPass {
		if err := g.CanPass(userID); err != nil {
			return err
		}

		p.Passed = true

		if g.AllAttackersPassed() {
		}

		return nil
	}

	if err := g.ValidateParameters(userID, pairIndex, card); err != nil {
		return err
	}

	cp := t.TableState[*pairIndex]

	// DEFENSE: Only the DEFENDER during GAME phase and CARD defends corresponding PAIR
	if action == ActionDefend {
		if err := g.CanDefend(userID, *card, cp); err != nil {
			return err
		}

		g.Turn.TableState[*pairIndex].DefenseCard = card

		p.hand = p.RemoveCard(*card)

		g.ResetPassingState()

		return nil
	}

	// ATTACK: EITHER phase is INITIAL and user is the initial attacker OR phase is GAME and user is an ATTACKER and CARD is within valid pool
	if action == ActionAttack {
		if err := g.CanAttack(userID, *card, *pairIndex); err != nil {
			return err
		}

		g.Turn.TableState[*pairIndex].AttackCard = *card

		p.hand = p.RemoveCard(*card)

		if t.Phase == TurnPhaseInitial {
			g.Turn.Phase = TurnPhaseGame
		}

		g.ResetPassingState()

		return nil
	}

	return Error{
		Code:    "INVALID_ACTION",
		Message: fmt.Sprintf("Action event provided invalid action: %v", action),
	}
}
