package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type PayBailDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *PayBailDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionPayBail
}

func (n *PayBailDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	if !cp.InJail {
		return ErrPlayerNotInJail
	}

	if !cp.Wallet.Deduct(50) {
		return ErrInsufficientFunds
	}

	cp.InJail = false
	cp.JailTurnsServed = 0
	g.PendingDecision = state.DecisionNone

	return n.ExecuteNext(g)
}
