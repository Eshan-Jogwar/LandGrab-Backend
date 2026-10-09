package movement

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type MoveManager struct {
	chain pipeline.MoveNode
}

// DefaultMoveManager is the shared, stateless flyweight MoveManager instance using DefaultDiceRoller.
var DefaultMoveManager = NewMoveManager(DefaultDiceRoller)

// Execute executes the movement chain on the given game using the default flyweight manager.
func Execute(g *state.Game) error {
	return DefaultMoveManager.Execute(g)
}

// NewMoveManager constructs the linear movement chain. An optional roller can be provided for deterministic tests.
func NewMoveManager(roller DiceRoller) *MoveManager {
	diceNode := NewDiceRollNode(roller)
	jailNode := &JailCheckNode{}
	posNode := &PositionUpdateNode{}
	goToJailNode := &GoToJailNode{}
	taxNode := &TaxNode{}
	propNode := &PropertyLandingNode{}
	railNode := &RailwayLandingNode{}
	utilNode := &UtilityLandingNode{}
	turnEndNode := &TurnEndNode{}

	chain := pipeline.BuildMoveChain(
		diceNode,
		jailNode,
		posNode,
		goToJailNode,
		taxNode,
		propNode,
		railNode,
		utilNode,
		turnEndNode,
	)

	return &MoveManager{chain: chain}
}

func (m *MoveManager) Execute(g *state.Game) error {
	if m.chain == nil {
		return nil
	}
	return m.chain.Execute(g)
}
