package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type TurnEndNode struct {
	pipeline.BaseMoveNode
}

func (n *TurnEndNode) ShouldExecute(g *state.Game) bool {
	// Only wrap up if we are not waiting for player decision
	return g.Phase != state.PhaseAwaitingDecision
}

func (n *TurnEndNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	if g.CanRollAgain {
		// Player rolled doubles and is eligible to roll again
		g.Phase = state.PhaseReadyToRoll
	} else {
		g.Phase = state.PhaseTurnEnded
		g.AdvanceTurn()
	}

	return n.ExecuteNext(g)
}
