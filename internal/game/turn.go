package game

type Turn struct {
	TableState        []CardPair `json:"table_state"`
	DefenderID        string     `json:"defender_id"`
	InitialAttackerID string     `json:"initial_attacker_id"`

	AttackerIDs map[string]bool `json:"attacker_ids"` // Figure out how to do this, i doubt this is a good pattern
	Phase       TurnPhase       `json:"phase"`
}

type TurnPhase int

const (
	TurnPhaseInitial  TurnPhase = iota // This represents the initial state of the round (initial attacker has yet to place first card, no one but initial attacker can attack)
	TurnPhaseGame                      // This represents that the game is active (i.e. defender and attackers can make simultaneous actions)
	TurnPhaseTake                      // This represents when a defender either opts to take or cannot defend the attacks placed down, gives a small period where attackers can toss on additional cards
	TurnPhaseDraw                      // All players that need to draw, draw up to 6 cards
	TurnPhaseComplete                  // Represents a completed combat, about to go to a new inital state
)

type Action int

const (
	ActionDefend Action = iota
	ActionAttack
	ActionTake
	ActionPass
)

// Checks if there is still capacity for another attack card
func (t *Turn) TableIsFull() bool {
	for _, pair := range t.TableState {
		if !pair.IsOccupied() {
			return false
		}
	}

	return true
}

func (t *Turn) GetPairFromIndex(index int) CardPair {
	return t.TableState[index]
}

func (t *Turn) GetPlayedCards() []Card {
	var cards []Card

	for _, pair := range t.TableState {
		if !pair.IsOccupied() {
			continue
		}

		cards = append(cards, pair.AttackCard)
		if pair.DefenseCard != nil && !pair.DefenseCard.isEmpty() {
			cards = append(cards, *pair.DefenseCard)
		}
	}

	return cards
}

func (g *Game) ResetPassingState() {
	for _, pl := range g.players {
		pl.Passed = false
	}
}

func (g *Game) AllAttackersPassed() bool {
	for id := range g.Turn.AttackerIDs {
		if !g.players[id].Passed {
			return false
		}
	}

	return true
}
