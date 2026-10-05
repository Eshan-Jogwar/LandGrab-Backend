package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/types"
)

// transferRent transfers rent between players, handling bankruptcy if funds are insufficient.
func transferRent(from *state.Player, to *state.Player, rent types.Money, g *state.Game) {
	if from.Wallet.Deduct(rent) {
		to.Wallet.Add(rent)
		return
	}
	// From player cannot pay full rent
	remaining := from.Wallet.Balance()
	from.Wallet.Deduct(remaining)
	to.Wallet.Add(remaining)
	from.IsBankrupt = true
	g.Phase = state.PhaseTurnEnded
}

type PropertyLandingNode struct {
	pipeline.BaseMoveNode
}

func (n *PropertyLandingNode) ShouldExecute(g *state.Game) bool {
	if g.Phase == state.PhaseTurnEnded {
		return false
	}
	t := g.CurrentTile()
	return t != nil && state.IsStreetProperty(t)
}

func (n *PropertyLandingNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	t := g.CurrentTile()
	owner, ph := g.GetPropertyOwner(t.ID())

	if owner == nil {
		// Unowned property: pause and prompt player decision
		g.Phase = state.PhaseAwaitingDecision
		g.PendingDecision = state.DecisionBuy
		g.DecisionTargetTile = t.ID()
		return n.ExecuteNext(g)
	}

	if owner.ID != cp.ID {
		// Opponent owned property: calculate and transfer rent
		rent, err := ph.CalculateRent(g.DiceRoll[0], g.DiceRoll[1], owner.CountRailways(), owner.CountUtilities())
		if err == nil && rent > 0 {
			transferRent(cp, owner, rent, g)
		}
	}

	return n.ExecuteNext(g)
}

type RailwayLandingNode struct {
	pipeline.BaseMoveNode
}

func (n *RailwayLandingNode) ShouldExecute(g *state.Game) bool {
	if g.Phase == state.PhaseTurnEnded {
		return false
	}
	t := g.CurrentTile()
	return t != nil && state.IsRailway(t)
}

func (n *RailwayLandingNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	t := g.CurrentTile()
	owner, ph := g.GetPropertyOwner(t.ID())

	if owner == nil {
		g.Phase = state.PhaseAwaitingDecision
		g.PendingDecision = state.DecisionBuy
		g.DecisionTargetTile = t.ID()
		return n.ExecuteNext(g)
	}

	if owner.ID != cp.ID {
		rent, err := ph.CalculateRent(g.DiceRoll[0], g.DiceRoll[1], owner.CountRailways(), 0)
		if err == nil && rent > 0 {
			transferRent(cp, owner, rent, g)
		}
	}

	return n.ExecuteNext(g)
}

type UtilityLandingNode struct {
	pipeline.BaseMoveNode
}

func (n *UtilityLandingNode) ShouldExecute(g *state.Game) bool {
	if g.Phase == state.PhaseTurnEnded {
		return false
	}
	t := g.CurrentTile()
	return t != nil && state.IsUtility(t)
}

func (n *UtilityLandingNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	t := g.CurrentTile()
	owner, ph := g.GetPropertyOwner(t.ID())

	if owner == nil {
		g.Phase = state.PhaseAwaitingDecision
		g.PendingDecision = state.DecisionBuy
		g.DecisionTargetTile = t.ID()
		return n.ExecuteNext(g)
	}

	if owner.ID != cp.ID {
		rent, err := ph.CalculateRent(g.DiceRoll[0], g.DiceRoll[1], 0, owner.CountUtilities())
		if err == nil && rent > 0 {
			transferRent(cp, owner, rent, g)
		}
	}

	return n.ExecuteNext(g)
}
