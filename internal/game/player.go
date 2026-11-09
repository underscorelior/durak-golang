package game

import (
	"cmp"
	"maps"
	"slices"
)

type Player struct {
	hand     []Card
	Position int
	Passed   bool
}

type PlayerView struct {
	PlayerID PlayerID `json:"player_id"`
	HandSize int      `json:"hand_size"`
	Position int      `json:"position"`
	Passed   bool     `json:"passed"`
}

type Snapshot struct {
	Players      []PlayerView `json:"players"`
	Hand         []Card       `json:"hand"`
	Trump        Card         `json:"trump"`
	DeckSize     int          `json:"deck_size"`
	Turn         Turn         `json:"turn"`
	FirstDefense bool         `json:"first_defense"`
}

func (g *Game) playerSnapshots() []PlayerView {
	var playerSnapshots []PlayerView
	for playerID, player := range g.players {
		playerSnapshot := PlayerView{
			PlayerID: playerID,
			HandSize: len(player.hand),
			Position: player.Position,
			Passed:   player.Passed,
		}
		playerSnapshots = append(playerSnapshots, playerSnapshot)
	}

	return playerSnapshots
}

func (g *Game) SnapshotFor(playerID PlayerID) *Snapshot { // TODO: Two versions of snapshot, lobby and game. need to fix naming.
	if _, ok := g.players[playerID]; !ok {
		return nil
	}

	return &Snapshot{
		Players:      g.playerSnapshots(),
		Hand:         g.players[playerID].hand,
		Trump:        g.Trump,
		DeckSize:     g.DeckSize(),
		Turn:         g.Turn,
		FirstDefense: g.FirstDefense,
	}
}

func (p *Player) RemoveCard(card Card) []Card {
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

func (p *Player) CardInHand(card Card) bool {
	return slices.Contains(p.hand, card)
}

type PlayerOrder struct {
	Next     PlayerID
	Previous PlayerID
}
type TableOrder map[PlayerID]PlayerOrder

func (g *Game) PlayerAtPosition(pos int) PlayerID {
	for id, player := range g.players {
		if player.Position == pos {
			return id
		}
	}

	return ""
}

func (g *Game) SetupTableOrder() {
	numPlayers := len(g.players)
	tableOrder := make(TableOrder, numPlayers)

	players := g.players

	keys := slices.Collect(maps.Keys(players))

	slices.SortFunc(keys, func(i, j PlayerID) int {
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

func (g *Game) AdjacentTo(playerID PlayerID) []PlayerID {
	o := g.Order[playerID]
	return []PlayerID{o.Next, o.Previous}
}

func calcPosOffset(pos, offset, num int) int {
	return ((pos+offset)%num + num) % num
}
