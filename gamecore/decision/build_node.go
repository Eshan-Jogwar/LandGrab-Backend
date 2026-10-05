package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type BuildHouseDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *BuildHouseDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionBuildHouse
}

func (n *BuildHouseDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	targetID := g.DecisionTargetTile
	targetPH, ok := cp.FindProperty(targetID)
	if !ok {
		return ErrPropertyNotOwned
	}

	if !targetPH.Tile.Buildable() {
		return ErrInvalidDecision
	}

	color := targetPH.Tile.Color()

	// 1. Monopoly Rule Check
	if cp.CountColorGroup(color) != state.TotalInColorGroup(color) {
		return ErrMonopolyRequired
	}

	// 2. No Mortgage in Color Group Rule Check
	for _, ph := range cp.Properties {
		if ph.Tile.Color() == color && ph.IsMortgaged {
			return ErrMortgagedInGroup
		}
	}

	// 3. Even Building Rule Check
	// Effective level: 0-4 houses, 5 for hotel
	effectiveLevel := func(ph *state.PropertyHolder) int {
		if ph.ExtraInfo.Hotel {
			return 5
		}
		return ph.ExtraInfo.Houses
	}

	minLevel := 5
	for _, ph := range cp.Properties {
		if ph.Tile.Color() == color {
			lvl := effectiveLevel(ph)
			if lvl < minLevel {
				minLevel = lvl
			}
		}
	}

	if effectiveLevel(targetPH) > minLevel {
		return ErrEvenBuildingRule
	}

	// 4. Capacity & Cost Deduction
	if targetPH.ExtraInfo.Hotel {
		return ErrMaxHousesReached
	}

	if targetPH.ExtraInfo.Houses == 4 {
		// Upgrading 4 houses to 1 Hotel
		hotelCost, err := targetPH.Tile.CostHotel()
		if err != nil {
			return err
		}
		if !cp.Wallet.Deduct(hotelCost) {
			return ErrInsufficientFunds
		}
		targetPH.ExtraInfo.Hotel = true
		targetPH.ExtraInfo.Houses = 0
	} else {
		// Adding 1 House
		houseCost, err := targetPH.Tile.CostHouse()
		if err != nil {
			return err
		}
		if !cp.Wallet.Deduct(houseCost) {
			return ErrInsufficientFunds
		}
		targetPH.ExtraInfo.Houses++
	}

	g.PendingDecision = state.DecisionNone
	return n.ExecuteNext(g)
}
