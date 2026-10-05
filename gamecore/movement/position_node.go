package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/tile"
)

type PositionUpdateNode struct {
	pipeline.BaseMoveNode
}

func (n *PositionUpdateNode) ShouldExecute(g *state.Game) bool {
	cp := g.CurrentPlayer()
	return cp != nil && !cp.InJail && g.Phase != state.PhaseTurnEnded
}

func (n *PositionUpdateNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	rollSum := int(g.DiceRoll[0] + g.DiceRoll[1])
	newPos := cp.Position + rollSum

	// If passed or landed on Go (circular board wrap), collect $200
	if newPos >= tile.BoardTileCount {
		cp.Wallet.Add(200)
	}

	cp.Position = newPos % tile.BoardTileCount
	return n.ExecuteNext(g)
}
