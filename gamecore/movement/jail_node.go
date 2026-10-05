package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type JailCheckNode struct {
	pipeline.BaseMoveNode
}

func (n *JailCheckNode) ShouldExecute(g *state.Game) bool {
	cp := g.CurrentPlayer()
	return cp != nil && cp.InJail && g.Phase != state.PhaseTurnEnded
}

func (n *JailCheckNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	isDoubles := g.DiceRoll[0] == g.DiceRoll[1]

	if isDoubles {
		// Released by rolling doubles
		cp.InJail = false
		cp.JailTurnsServed = 0
		g.ConsecutiveDoubles = 0
		g.CanRollAgain = false // Rolling doubles to exit jail does not grant an extra turn
		return n.ExecuteNext(g)
	}

	// Did not roll doubles
	cp.JailTurnsServed++
	if cp.JailTurnsServed >= 3 {
		// Must pay $50 bail on the 3rd turn
		if cp.Wallet.Deduct(50) {
			cp.InJail = false
			cp.JailTurnsServed = 0
			// Allowed to move with current roll
			return n.ExecuteNext(g)
		}
		// Insolvent/Bankrupt if cannot pay $50 bail
		cp.IsBankrupt = true
		g.Phase = state.PhaseTurnEnded
		return n.ExecuteNext(g)
	}

	// Still in jail, turn ends without moving
	g.Phase = state.PhaseTurnEnded
	return n.ExecuteNext(g)
}
