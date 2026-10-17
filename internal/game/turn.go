package game

import (
	"maps"
	"slices"
	"sort"
)

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
	ActionDefense Action = iota
	ActionAttack
	ActionTake
	ActionPass
)

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

// TODO: Better name
type PlayerOrder struct {
	Next     string
	Previous string
}
type TableOrder map[string]PlayerOrder

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

	sort.Slice(keys[:], func(i, j int) bool {
		return players[keys[i]].Position < players[keys[j]].Position
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
