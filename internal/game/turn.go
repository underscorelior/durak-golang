package game

type Turn struct {
	TableState        []CardPair
	DefenderID        string // Is this the best pattern?
	InitialAttackerID string

	AttackerIDs map[string]struct{} // Figure out how to do this, i doubt this is a good pattern
	Phase       TurnPhase
}

type TurnPhase uint8

const (
	INITIAL = iota
	ATTACK
	DEFENSE
	COMPLETE
)

// Checks if the player is a valid attacker
func (t *Turn) IsAttacker(player Player) bool {
	_, exists := t.AttackerIDs[player.ID]
	return exists
}

// Checks if the player is the current defender
func (t *Turn) IsDefender(player Player) bool {
	return player.ID == t.InitialAttackerID
}
