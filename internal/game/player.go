package game

type Player struct {
	ID       string
	Position uint8 // Need a different pattern

	hand []Card
}
