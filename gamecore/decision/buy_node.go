package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/tile"
)

type BuyDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *BuyDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionBuy
}

func (n *BuyDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	targetTile := tile.TileAt(int(g.DecisionTargetTile))
	cost, err := targetTile.Cost()
	if err != nil {
		return err
	}

	if !cp.Wallet.Deduct(cost) {
		return ErrInsufficientFunds
	}

	cp.AddProperty(targetTile)
	g.TileOwnerIndex[targetTile.ID()] = cp.ID
	g.PendingDecision = state.DecisionNone

	if !g.CanRollAgain {
		g.Phase = state.PhaseTurnEnded
		g.AdvanceTurn()
	} else {
		g.Phase = state.PhaseReadyToRoll
	}

	return n.ExecuteNext(g)
}

type DeclineBuyDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *DeclineBuyDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionDeclineBuy
}

func (n *DeclineBuyDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	g.PendingDecision = state.DecisionNone

	if !g.CanRollAgain {
		g.Phase = state.PhaseTurnEnded
		g.AdvanceTurn()
	} else {
		g.Phase = state.PhaseReadyToRoll
	}

	return n.ExecuteNext(g)
}
