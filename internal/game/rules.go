package game

// Can the defense Card beat the attack Card, checks Trump supremacy too
func (g *Game) CanBeat(attack Card, defense Card) bool {
	// If the suits are equal, return whether the defense is more powerful than the attack.
	if attack.Suit == defense.Suit {
		return attack.Rank < defense.Rank
	}

	// Else return whether the defense is a trump card
	return defense.Suit == g.Trump.Suit
}

// Checks if a player is allowed to place (Both existing ranks and if its attacker)
func (t *Turn) CanPlace(player Player, card Card) bool {
	if t.Phase == INITIAL && player.ID == t.InitialAttackerID {
		return true
	}

	rankSet := make(map[Rank]struct{}) // Should this be bool isntead of struct{}?
	for _, pairs := range t.TableState {
		rankSet[pairs.AttackCard.Rank] = struct{}{}
		rankSet[pairs.DefenseCard.Rank] = struct{}{}
	}

	_, exists := rankSet[card.Rank]

	// TODO: This would mean that no one can attack while there is a card that has not been defended,
	// 		 need to find this out. Also might cause an issue with multiple cards.
	return exists && t.IsAttacker(player) && t.Phase == ATTACK
}

// Checks if the Player is allowed to make an attack
//   - Checks if the player is a defender
//   - Checks if the player is the initial attacker and the current phase is INITIAL
//   - Checks if the player is an attacker and if the Rank exists
func (g *Game) CanAttack(player Player, attack Card) bool {
	if g.Turn.IsDefender(player) {
		return false
	}

	return g.Turn.CanPlace(player, attack)
}

// Checks if the player is capable of making a valid defense
func (g *Game) CanDefend(player Player, attack Card, defense Card) bool {
	if !g.Turn.IsDefender(player) {
		return false
	}

	return g.CanBeat(attack, defense) && g.Turn.Phase == DEFENSE
}

// Checks if a player is allowed to take
func (g *Game) CanTake(player Player) bool {
	if !g.Turn.IsDefender(player) {
		return false
	}

	return g.Turn.Phase == DEFENSE
}
