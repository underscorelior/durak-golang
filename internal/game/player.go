package game

import (
	"cmp"
	"maps"
	"slices"
)

type PlayerState struct {
	hand     []Card
	Position int
	Passed   bool
}

type GamePlayerState struct {
	UserID   string `json:"user_id"`
	HandSize int    `json:"hand_size"`
	Position int    `json:"position"`
	Passed   bool   `json:"passed"`
}

type GameStateSnapshot struct {
	Players      []GamePlayerState `json:"players"`
	Hand         []Card            `json:"hand"`
	Trump        Card              `json:"trump"`
	DeckSize     int               `json:"deck_size"`
	Turn         Turn              `json:"turn"`
	FirstDefense bool              `json:"first_defense"`
}

func (g *Game) playerSnapshots() []GamePlayerState {
	var playerSnapshots []GamePlayerState
	for userID, player := range g.players {
		playerSnapshot := GamePlayerState{
			UserID:   userID,
			HandSize: len(player.hand),
			Position: player.Position,
			Passed:   player.Passed,
		}
		playerSnapshots = append(playerSnapshots, playerSnapshot)
	}

	return playerSnapshots
}

func (g *Game) StateFor(userID string) *GameStateSnapshot {
	_, ok := g.players[userID]

	if !ok {
		// TODO: THROW ERROR HERE!!!
		return &GameStateSnapshot{}
	}

	return &GameStateSnapshot{
		Players:      g.playerSnapshots(),
		Hand:         g.players[userID].hand,
		Trump:        g.Trump,
		DeckSize:     g.DeckSize(),
		Turn:         g.Turn,
		FirstDefense: g.FirstDefense,
	}
}

func (p *PlayerState) RemoveCard(card Card) []Card {
	hand := slices.Collect(func(yield func(Card) bool) {
		for _, c := range p.hand {
			if c.Rank != card.Rank || c.Suit != card.Suit {
				if !yield(c) {
					return
				}
			}
		}
	})

	return hand
}

func (p *PlayerState) CardInHand(card Card) bool {
	return slices.Contains(p.hand, card)
}

type PlayerOrder struct {
	Next     string
	Previous string
}
type TableOrder map[string]PlayerOrder

// TODO: Fix this to be better O(n) wise.
func (g *Game) PlayerAtPosition(pos int) string {
	for id, player := range g.players {
		if player.Position == pos {
			return id
		}
	}

	return ""
}

// TODO: Account for gaps in position
func (g *Game) SetupTableOrder() {
	numPlayers := len(g.players)
	tableOrder := make(TableOrder, numPlayers)

	players := g.players

	keys := slices.Collect(maps.Keys(players))

	slices.SortFunc(keys, func(i, j string) int {
		return cmp.Compare(players[i].Position, players[j].Position)
	})

	for _, key := range keys {
		pos := players[key].Position
		tableOrder[key] = PlayerOrder{
			Next:     g.PlayerAtPosition(calcPosOffset(pos, 1, numPlayers)),
			Previous: g.PlayerAtPosition(calcPosOffset(pos, -1, numPlayers)),
		}
	}

	g.Order = tableOrder
}

func (g *Game) AdjacentTo(userID string) map[string]bool {
	o := g.Order[userID]
	return map[string]bool{o.Next: true, o.Previous: true}
}

func calcPosOffset(pos, offset, num int) int {
	return ((pos+offset)%num + num) % num
}
