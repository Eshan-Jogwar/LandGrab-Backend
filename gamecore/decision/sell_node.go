package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type SellHouseDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *SellHouseDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionSellHouse
}

func (n *SellHouseDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	cp := g.CurrentPlayer()
	targetID := g.DecisionTargetTile
	targetPH, ok := cp.FindProperty(targetID)
	if !ok {
		return ErrPropertyNotOwned
	}

	effectiveLevel := func(ph *state.PropertyHolder) int {
		if ph.ExtraInfo.Hotel {
			return 5
		}
		return ph.ExtraInfo.Houses
	}

	if effectiveLevel(targetPH) == 0 {
		return ErrNoHousesToSell
	}

	color := targetPH.Tile.Color()

	// Even Selling Rule: must sell from the property with the maximum houses in the group
	maxLevel := 0
	for _, ph := range cp.Properties {
		if ph.Tile.Color() == color {
			lvl := effectiveLevel(ph)
			if lvl > maxLevel {
				maxLevel = lvl
			}
		}
	}

	if effectiveLevel(targetPH) < maxLevel {
		return ErrEvenSellingRule
	}

	if targetPH.ExtraInfo.Hotel {
		// Downgrade hotel to 4 houses, refund 50% hotel cost
		hotelCost, _ := targetPH.Tile.CostHotel()
		targetPH.ExtraInfo.Hotel = false
		targetPH.ExtraInfo.Houses = 4
		cp.Wallet.Add(hotelCost / 2)
	} else {
		// Sell 1 house, refund 50% house cost
		houseCost, _ := targetPH.Tile.CostHouse()
		targetPH.ExtraInfo.Houses--
		cp.Wallet.Add(houseCost / 2)
	}

	g.PendingDecision = state.DecisionNone
	return n.ExecuteNext(g)
}
