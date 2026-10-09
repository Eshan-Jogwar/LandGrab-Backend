package state

import (
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

type TradeOffer struct {
	FromPlayerID types.ID
	ToPlayerID   types.ID
	TileID       types.ID
	Cash         types.Money
}

type Game struct {
	Players            []*Player
	CurrentPlayerIdx   int
	Phase              GamePhase
	PendingDecision    DecisionType
	DecisionTargetTile types.ID
	PendingTrade       *TradeOffer

	// TileOwnerIndex maps tile ID [0..39] to Owner Player ID or NoOwnerID.
	TileOwnerIndex [tile.BoardTileCount]types.ID

	// Turn movement state
	DiceRoll           [2]types.Dice
	ConsecutiveDoubles int
	CanRollAgain       bool
}

// NewGame initializes a new game with the given players and resets ownership to NoOwnerID.
func NewGame(players []*Player) *Game {
	g := &Game{
		Players:          players,
		CurrentPlayerIdx: 0,
		Phase:            PhaseReadyToRoll,
		PendingDecision:  DecisionNone,
	}
	for i := range g.TileOwnerIndex {
		g.TileOwnerIndex[i] = NoOwnerID
	}
	return g
}

func (g *Game) CurrentPlayer() *Player {
	if len(g.Players) == 0 {
		return nil
	}
	return g.Players[g.CurrentPlayerIdx]
}

func (g *Game) CurrentTile() tile.Tile {
	cp := g.CurrentPlayer()
	if cp == nil {
		return nil
	}
	return tile.TileAt(cp.Position)
}

func (g *Game) GetPropertyOwner(tileID types.ID) (*Player, *PropertyHolder) {
	if int(tileID) >= tile.BoardTileCount {
		return nil, nil
	}
	ownerID := g.TileOwnerIndex[tileID]
	if ownerID == NoOwnerID {
		return nil, nil
	}
	for _, p := range g.Players {
		if p.ID == ownerID {
			ph, _ := p.FindProperty(tileID)
			return p, ph
		}
	}
	return nil, nil
}

func (g *Game) AdvanceTurn() {
	if len(g.Players) == 0 {
		return
	}
	g.CurrentPlayerIdx = (g.CurrentPlayerIdx + 1) % len(g.Players)
	g.Phase = PhaseReadyToRoll
	g.ConsecutiveDoubles = 0
	g.CanRollAgain = false
	g.PendingDecision = DecisionNone
	g.DecisionTargetTile = 0
	g.PendingTrade = nil
}
