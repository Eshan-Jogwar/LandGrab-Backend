package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type MortgageDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *MortgageDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionMortgage
}

func (n *MortgageDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	targetID := g.DecisionTargetTile
	targetPH, ok := cp.FindProperty(targetID)
	if !ok {
		return ErrPropertyNotOwned
	}

	if targetPH.IsMortgaged {
		return ErrAlreadyMortgaged
	}

	// Cannot mortgage if ANY property in the color group has buildings
	color := targetPH.Tile.Color()
	if color != 0 {
		for _, ph := range cp.Properties {
			if ph.Tile.Color() == color {
				if ph.ExtraInfo.Houses > 0 || ph.ExtraInfo.Hotel {
					return ErrBuildingsInGroup
				}
			}
		}
	}

	mortVal, err := targetPH.Tile.MortgageValue()
	if err != nil {
		return err
	}

	targetPH.IsMortgaged = true
	cp.Wallet.Add(mortVal)
	g.PendingDecision = state.DecisionNone

	return n.ExecuteNext(g)
}

type DemortgageDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *DemortgageDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionDemortgage
}

func (n *DemortgageDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	targetID := g.DecisionTargetTile
	targetPH, ok := cp.FindProperty(targetID)
	if !ok {
		return ErrPropertyNotOwned
	}

	if !targetPH.IsMortgaged {
		return ErrNotMortgaged
	}

	demortVal, err := targetPH.Tile.DemortgageValue()
	if err != nil {
		return err
	}

	if !cp.Wallet.Deduct(demortVal) {
		return ErrInsufficientFunds
	}

	targetPH.IsMortgaged = false
	g.PendingDecision = state.DecisionNone

	return n.ExecuteNext(g)
}
