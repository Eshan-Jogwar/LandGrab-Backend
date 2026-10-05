package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/types"
)

type GoToJailNode struct {
	pipeline.BaseMoveNode
}

func (n *GoToJailNode) ShouldExecute(g *state.Game) bool {
	cp := g.CurrentPlayer()
	return cp != nil && cp.Position == 30 && g.Phase != state.PhaseTurnEnded
}

func (n *GoToJailNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	cp.Position = 10
	cp.InJail = true
	cp.JailTurnsServed = 0
	g.ConsecutiveDoubles = 0
	g.CanRollAgain = false
	g.Phase = state.PhaseTurnEnded
	return n.ExecuteNext(g)
}

type TaxNode struct {
	pipeline.BaseMoveNode
}

func (n *TaxNode) ShouldExecute(g *state.Game) bool {
	cp := g.CurrentPlayer()
	if cp == nil || g.Phase == state.PhaseTurnEnded {
		return false
	}
	return cp.Position == 4 || cp.Position == 38
}

func (n *TaxNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	tax := types.Money(200)
	if cp.Position == 38 {
		tax = 100 // Wealth Tax
	}

	if !cp.Wallet.Deduct(tax) {
		cp.IsBankrupt = true
		g.Phase = state.PhaseTurnEnded
		return n.ExecuteNext(g)
	}

	return n.ExecuteNext(g)
}
