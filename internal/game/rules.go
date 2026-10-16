package game

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

func (g *Game) CanTake(userID string) (bool, Error) {
	t := g.Turn

	if userID != t.DefenderID {
		return false, Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not the defender",
		}
	}

	if t.Phase != TP_GAME {
		return false, Error{
			Code:    "INVALID_PHASE",
			Message: "Cannot take cards when game is not real-time",
		}
	}

	return true, Error{}
}

func (g *Game) CanPass(userID string) (bool, Error) {
	t := g.Turn

	if _, exists := t.AttackerIDs[userID]; !exists {
		return false, Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not an attacker",
		}
	}

	if t.Phase != TP_GAME {
		return false, Error{
			Code:    "INVALID_PHASE",
			Message: "Cannot pass turn when game is not real-time",
		}
	}

	return true, Error{}
}

func (g *Game) CanDefend(userID string, card Card, cardPair CardPair) (bool, Error) {
	t := g.Turn

	if userID != t.DefenderID {
		return false, Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not the defender",
		}
	}

	if t.Phase != TP_GAME {
		return false, Error{
			Code:    "INVALID_PHASE",
			Message: "Cannot defend when game is not real-time",
		}
	}

	if cardPair.DefenseCard != nil {
		return false, Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair has already been defended",
		}
	}

	if !cardPair.DefendedBy(card, g.Trump.Suit) {
		return false, Error{
			Code:    "INVALID_DEFENSE_CARD",
			Message: "Provided defense card is unable to beat attack card",
		}
	}

	return true, Error{}
}

func (g *Game) CanAttack(userID string, card Card, pairIndex int) (bool, Error) {
	t := g.Turn

	if t.Phase == TP_INITIAL && userID == t.InitialAttackerID {
		return true, Error{}
	}

	if _, exists := t.AttackerIDs[userID]; !exists {
		return false, Error{
			Code:    "INVALID_ROLE",
			Message: "This user is not an attacker",
		}
	}

	if t.Phase != TP_GAME {
		return false, Error{
			Code:    "INVALID_PHASE",
			Message: "Cannot attack when game is not real-time",
		}
	}

	if t.TableIsFull() {
		return false, Error{
			Code:    "TABLE_FULL",
			Message: "Unable to place additional attack card",
		}
	}

	if !t.CanPlace(card) {
		return false, Error{
			Code:    "INVALID_CARD",
			Message: "This card's rank is not placeable in the current conditions",
		}
	}

	if g.Turn.TableState[pairIndex].IsOccupied() {
		return false, Error{
			Code:    "OCCUPIED_PAIR",
			Message: "This card pair already contains an attack card",
		}
	}

	return true, Error{}
}

func (g *Game) MakeAction(userID string, action Action, card *Card, pairIndex *int) (bool, Error) {
	t := g.Turn
	p := g.players[userID]

	if p == nil {
		return false, Error{
			Code:    "INVALID_USER",
			Message: "Provided user ID is not in the game",
		}
	}

	// Take and Pass

	// TAKE: Only DEFENDER during GAME phase
	if action == A_TAKE {
		ok, err := g.CanTake(userID)

		if !ok {
			return false, err
		}

		// START TIMER ALLOWING ATTACKERS TO THROW IN
		// NEED A WAY TO ACCURATELY DISPLAY IT (SENDING TIMESTAMP AND ALLOWING CLIENT TO CALC HOW LONG IS LEFT?)
		// ONCE TIME IS UP, MOVE CARDS TO PLAYER HAND
		// p.hand = append(p.hand, t.GetPlayedCards()...)

		return true, Error{}
	}

	// PASS: Only ATTACKERs during GAME phase, PASS BY DEFAULT ONCE CANNOT ADD NEW CARDS
	if action == A_PASS {
		ok, err := g.CanPass(userID)

		if !ok {
			return false, err
		}

		p.Passed = true

		return true, Error{}
	}

	// Defense and Attack
	if pairIndex == nil {
		return false, Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing CardPair index",
		}
	}
	if card == nil {
		return false, Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event missing Card played",
		}
	}

	if !p.CardInHand(*card) {
		return false, Error{
			Code:    "CARD_NOT_IN_HAND",
			Message: "Selected card not in players hand",
		}
	}

	if *pairIndex < 0 || *pairIndex > 5 {
		return false, Error{
			Code:    "INVALID_RANGE",
			Message: "Provided index not within range",
		}
	}

	cp := t.TableState[*pairIndex]

	// DEFENSE: Only the DEFENDER during GAME phase and CARD defends corresponding PAIR
	if action == A_DEFENSE {
		ok, err := g.CanDefend(userID, *card, cp)

		if !ok {
			return false, err
		}

		g.Turn.TableState[*pairIndex].DefenseCard = card

		p.hand = p.RemoveCard(*card)

		return true, Error{}
	}

	// ATTACK: EITHER phase is INITIAL and user is the initial attacker OR phase is GAME and user is an ATTACKER and CARD is within valid pool
	if action == A_ATTACK {
		ok, err := g.CanAttack(userID, *card, *pairIndex)

		if !ok {
			return false, err
		}

		g.Turn.TableState[*pairIndex].AttackCard = *card

		p.hand = p.RemoveCard(*card)

		return true, Error{}
	}

	return false, Error{
		Code:    "INVALID_ACTION",
		Message: "Action event provided invalid action",
	}
}
