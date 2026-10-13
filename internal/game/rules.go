package game

// Checks if a player is allowed to be placed among existing cards
func (t *Turn) CanPlace(card Card) bool {
	rankSet := make(map[Rank]bool)

	for _, pairs := range t.TableState {
		rankSet[pairs.AttackCard.Rank] = true
		rankSet[pairs.DefenseCard.Rank] = true
	}

	_, exists := rankSet[card.Rank]

	return exists
}

func (g *Game) ValidateTurn(userID string, action Action, card *Card, pairIndex *int) (bool, Error) {
	t := g.Turn

	// Take and Pass
	if action == Take {
		if !t.IsDefender(userID) {
			// Error message for this
			return false, Error{
				Code:    "INVALID_ROLE",
				Message: "This user is not the defender",
			}
		}

		if t.Phase != DEFENSE {
			return false, Error{
				Code:    "INVALID_PHASE",
				Message: "Cannot take cards when it's not the defender's action",
			}
		}

		return true, Error{}
	}

	// TODO: Passing

	// Defense and Attack
	var cardPair CardPair

	if pairIndex == nil {
		return false, Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event is missing CardPair index",
		}
	}
	if card == nil {
		return false, Error{
			Code:    "MISSING_ARGUMENT",
			Message: "Action event is missing Card played",
		}
	}
	if g.CardInHand(userID, *card) {
		return false, Error{
			Code:    "CARD_NOT_IN_HAND",
			Message: "Selected card is not in players hand",
		}
	}
	cardPair = t.TableState[*pairIndex]

	// Defense:
	if action == Defense {
		if !t.IsDefender(userID) {
			return false, Error{
				Code:    "INVALID_ROLE",
				Message: "This user is not the defender",
			}
		}

		if t.Phase != DEFENSE {
			return false, Error{
				Code:    "INVALID_PHASE",
				Message: "Cannot defend when it's not the defender's action",
			}
		}

		if !cardPair.DefendedBy(*card, g.Trump.Suit) {
			return false, Error{
				Code:    "INVALID_DEFENSE_CARD",
				Message: "The provided defense card is unable to beat the attacking card",
			}
		}

		return true, Error{}
	}

	// Attack:
	if action == Attack {
		if t.Phase == INITIAL && userID == t.InitialAttackerID {
			return true, Error{}
		}

		if !t.IsAttacker(userID) {
			return false, Error{
				Code:    "INVALID_ROLE",
				Message: "This user is not an attacker",
			}
		}

		if t.Phase != ATTACK {
			return false, Error{
				Code:    "INVALID_PHASE",
				Message: "Cannot attack when it's not the attacker's action",
			}
		}

		if t.TableIsFull() {
			return false, Error{
				Code:    "TABLE_FULL",
				Message: "Unable to place additional attack card",
			}
		}

		if !t.CanPlace(*card) {
			return false, Error{
				Code:    "INVALID_CARD",
				Message: "This card's rank is not placeable in the current conditions",
			}
		}

		return true, Error{}
	}

	return false, Error{
		Code:    "INVALID_ACTION",
		Message: "Action event provided invalid action",
	}
}
