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
	INITIAL = iota
	ATTACK
	DEFENSE
	COMPLETE
)

type Action int

const (
	Defense = iota
	Attack
	Pass
	Take
)

// Checks if the player is a valid attacker
func (t *Turn) IsAttacker(playerID string) bool {
	_, exists := t.AttackerIDs[playerID]
	return exists
}

// Checks if the player is the current defender
func (t *Turn) IsDefender(playerID string) bool {
	return playerID == t.DefenderID
}

// Checks if there is still capacity for another attack card
func (t *Turn) TableIsFull() bool {
	// TODO: Support 5 before first beat and 6 after

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
	cards := make([]Card, 12)

	for _, pair := range t.TableState {
		cards = append(cards, pair.AttackCard, *pair.DefenseCard)
	}

	return cards
}
