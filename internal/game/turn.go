package game

type Turn struct {
	TableState        []CardPair `json:"table_state"`
	DefenderID        string     `json:"defender_id"` // Is this the best pattern?
	InitialAttackerID string     `json:"initial_attacker_id"`

	AttackerIDs map[string]struct{} // Figure out how to do this, i doubt this is a good pattern
	Phase       TurnPhase           `json:"phase"`
}

type TurnPhase int

const (
	INITIAL = iota
	ATTACK
	DEFENSE
	COMPLETE
)

// Checks if the player is a valid attacker
func (t *Turn) IsAttacker(playerID string) bool {
	_, exists := t.AttackerIDs[playerID]
	return exists
}

// Checks if the player is the current defender
func (t *Turn) IsDefender(playerID string) bool {
	return playerID == t.InitialAttackerID
}
