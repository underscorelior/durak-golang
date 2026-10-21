package game

import "fmt"

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

func (g *Game) CanTake(userID string) error {
	t := g.Turn

	if userID != t.DefenderID {
		return Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not the defender",
		}
	}

	if t.Phase != TurnPhaseGame {
		return Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot take cards when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	return nil
}

func (g *Game) CanPass(userID string) error {
	t := g.Turn

	if _, exists := t.AttackerIDs[userID]; !exists {
		return Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not an attacker",
		}
	}

	if t.Phase != TurnPhaseGame {
		return Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot pass turn when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if g.players[userID].Passed == true {
		return Error{
			Code:    "ALREADY_PASSED",
			Message: "Cannot pass turn while already passed", // TODO: Reword!
		}
	}

	return nil
}

func (g *Game) CanDefend(userID string, card Card, cardPair CardPair) error {
	t := g.Turn

	if userID != t.DefenderID {
		return Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not the defender",
		}
	}

	if t.Phase != TurnPhaseGame {
		return Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot defend when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if cardPair.DefenseCard != nil {
		return Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair has already been defended",
		}
	}

	if !cardPair.DefendedBy(card, g.Trump.Suit) {
		return Error{
			Code:    "INVALID_DEFENSE_CARD",
			Message: "Provided defense card is unable to beat attack card",
		}
	}

	return nil
}

func (g *Game) CanAttack(userID string, card Card, pairIndex int) error {
	t := g.Turn

	if t.Phase == TurnPhaseInitial && userID == t.InitialAttackerID {
		return nil
	}

	if _, exists := t.AttackerIDs[userID]; !exists {
		return Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not an attacker",
		}
	}

	if t.Phase != TurnPhaseGame {
		return Error{
			Code:    "INVALID_PHASE",
			Message: fmt.Sprintf("Cannot attack when game is not real-time. Current phase: %v", t.Phase),
		}
	}

	if t.TableIsFull() {
		return Error{
			Code:    "TABLE_FULL",
			Message: "Unable to place additional attack card",
		}
	}

	if !t.CanPlace(card) {
		return Error{
			Code:    "INVALID_CARD",
			Message: "This card's rank is not placeable in the current conditions",
		}
	}

	if g.Turn.TableState[pairIndex].IsOccupied() {
		return Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair already contains an attack card",
		}
	}

	return nil
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

		return nil
	}

	if pairIndex == nil {
		return Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing CardPair index",
		}
	}
	if card == nil {
		return Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing Card played",
		}
	}

	if !p.CardInHand(*card) {
		return Error{
			Code:    "CARD_NOT_IN_HAND",
			Message: "Selected card not in players hand",
		}
	}

	if *pairIndex < 0 || *pairIndex > 5 {
		return Error{
			Code:    "INVALID_RANGE",
			Message: fmt.Sprintf("Provided index (%v) not within range", *pairIndex),
		}
	}

	cp := t.TableState[*pairIndex]

	// DEFENSE: Only the DEFENDER during GAME phase and CARD defends corresponding PAIR
	if action == ActionDefense {
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
