package game

import (
	"fmt"
	"slices"
)

// Checks if a player is allowed to be placed among existing cards
func (t *Turn) CanPlace(card Card) bool {
	rankSet := make(map[Rank]bool)

	for _, pairs := range t.TableState {
		rankSet[pairs.AttackCard.Rank] = true
		if pairs.DefenseCard != nil {
			rankSet[pairs.DefenseCard.Rank] = true
		}
	}

	_, exists := rankSet[card.Rank]

	return exists
}

func (g *Game) CanTake(playerID PlayerID) error {
	t := g.Turn

	if playerID != t.DefenderID {
		return &Error{
			Code:    "INVALID_ROLE",
			Message: "This player is not the defender",
		}
	}

	if t.Phase != TurnPhaseGame {
		return &Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot take cards when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	return nil
}

func (g *Game) CanPass(playerID PlayerID) error {
	t := g.Turn

	if !slices.Contains(t.AttackerIDs, playerID) {
		return &Error{
			Code:    "INVALID_ROLE",
			Message: "This player is not an attacker",
		}
	}

	if t.Phase != TurnPhaseGame {
		return &Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot pass turn when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if g.players[playerID].Passed == true {
		return &Error{
			Code:    "ALREADY_PASSED",
			Message: "Cannot pass turn while already passed",
		}
	}

	return nil
}

func (g *Game) CanDefend(playerID PlayerID, card Card, cardPair CardPair) error {
	t := g.Turn

	if playerID != t.DefenderID {
		return &Error{
			Code:    "INVALID_ROLE",
			Message: "This player is not the defender",
		}
	}

	if t.Phase != TurnPhaseGame {
		return &Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot defend when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if cardPair.DefenseCard != nil {
		return &Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair has already been defended",
		}
	}

	if !cardPair.DefendedBy(card, g.Trump.Suit) {
		return &Error{
			Code:    "INVALID_DEFENSE_CARD",
			Message: "Provided defense card is unable to beat attack card",
		}
	}

	return nil
}

func (g *Game) CanAttack(playerID PlayerID, card Card, pairIndex int) error {
	t := g.Turn

	if t.Phase == TurnPhaseInitial && playerID == t.InitialAttackerID {
		return nil
	}

	if !slices.Contains(t.AttackerIDs, playerID) {
		return &Error{
			Code:    "INVALID_ROLE",
			Message: "This player is not an attacker",
		}
	}

	if t.Phase != TurnPhaseGame {
		return &Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot attack when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if t.TableIsFull() {
		return &Error{
			Code:    "TABLE_FULL",
			Message: "Unable to place additional attack card",
		}
	}

	if !t.CanPlace(card) {
		return &Error{
			Code:    "INVALID_CARD",
			Message: "This card's rank is not placeable in the current conditions",
		}
	}

	if g.Turn.TableState[pairIndex].IsOccupied() {
		return &Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair already contains an attack card",
		}
	}

	return nil
}

func (g *Game) ValidateParameters(playerID PlayerID, pairIndex *int, card *Card) error {
	if pairIndex == nil {
		return &Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing CardPair index",
		}
	}
	if card == nil {
		return &Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing Card played",
		}
	}

	if !g.players[playerID].CardInHand(*card) {
		return &Error{
			Code:    "CARD_NOT_IN_HAND",
			Message: "Selected card not in players hand",
		}
	}

	if *pairIndex < 0 || *pairIndex > 5 {
		return &Error{
			Code:    "INVALID_RANGE",
			Message: fmt.Sprintf("Provided index (%v) not within range", *pairIndex),
		}
	}

	return nil
}
