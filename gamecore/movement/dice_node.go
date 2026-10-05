package movement

import (
	"math/rand"

	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/types"
)

// DiceRoller is a function that generates two 6-sided dice rolls (1-6).
type DiceRoller func() (types.Dice, types.Dice)

func DefaultDiceRoller() (types.Dice, types.Dice) {
	d1 := types.Dice(rand.Intn(6) + 1)
	d2 := types.Dice(rand.Intn(6) + 1)
	return d1, d2
}

type DiceRollNode struct {
	pipeline.BaseMoveNode
	roller DiceRoller
}

func NewDiceRollNode(roller DiceRoller) *DiceRollNode {
	if roller == nil {
		roller = DefaultDiceRoller
	}
	return &DiceRollNode{roller: roller}
}

func (n *DiceRollNode) ShouldExecute(g *state.Game) bool {
	return g.Phase == state.PhaseReadyToRoll
}

func (n *DiceRollNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	d1, d2 := n.roller()
	g.DiceRoll = [2]types.Dice{d1, d2}

	isDoubles := d1 == d2
	if isDoubles {
		g.ConsecutiveDoubles++
		if g.ConsecutiveDoubles == 3 {
			// 3 consecutive doubles sends player directly to prison
			cp := g.CurrentPlayer()
			if cp != nil {
				cp.Position = 10
				cp.InJail = true
				cp.JailTurnsServed = 0
			}
			g.ConsecutiveDoubles = 0
			g.CanRollAgain = false
			g.Phase = state.PhaseTurnEnded
			return n.ExecuteNext(g)
		}
		g.CanRollAgain = true
	} else {
		g.ConsecutiveDoubles = 0
		g.CanRollAgain = false
	}

	return n.ExecuteNext(g)
}
