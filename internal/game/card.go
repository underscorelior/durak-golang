package game

type Suit uint8

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

type Rank uint8

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
	Suit Suit
	Rank Rank
}

type CardPair struct {
	AttackCard  Card
	DefenseCard *Card
	IsDefended  bool
}
