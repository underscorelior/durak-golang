package game

type Suit int

const (
	Club Suit = iota
	Diamond
	Heart
	Spade
)

var Suits = [4]Suit{
	Club,
	Diamond,
	Heart,
	Spade,
}

type Rank int

const (
	Six = iota + 6
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var Ranks = [9]Rank{
	Six,
	Seven,
	Eight,
	Nine,
	Ten,
	Jack,
	Queen,
	King,
	Ace,
}

type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// TODO: Dfiferentiate card by user who played it
type CardPair struct {
	AttackCard  Card  `json:"attack_card,omitempty"`
	DefenseCard *Card `json:"defense_card,omitempty"`
	IsDefended  bool  `json:"is_defended,omitempty"`
}

// Can the defense Card beat the attack Card, checks Trump supremacy too
func (cp *CardPair) DefendedBy(defense Card, trumpSuit Suit) bool {
	if cp.IsDefended {
		return false
	}

	attack := cp.AttackCard

	// If the suits are equal, return whether the defense is more powerful than the attack.
	if attack.Suit == defense.Suit {
		return attack.Rank < defense.Rank
	}

	// Else return whether the defense is a trump card
	return defense.Suit == trumpSuit
}

func (cp *CardPair) IsOccupied() bool {
	return cp.AttackCard.Rank != 0 && cp.AttackCard.Suit != 0 // FIXME: Need a better check than this.
}
