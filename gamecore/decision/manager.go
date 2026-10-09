package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type DecisionManager struct {
	chain pipeline.DecisionNode
}

// DefaultDecisionManager is the shared, stateless flyweight DecisionManager instance.
var DefaultDecisionManager = NewDecisionManager()

// Execute executes the decision chain on the given game using the default flyweight manager.
func Execute(g *state.Game) error {
	return DefaultDecisionManager.Execute(g)
}

// ExecuteTrade executes a trade decision on the given game using the default flyweight manager.
func ExecuteTrade(g *state.Game, trade *TradeOffer) error {
	return DefaultDecisionManager.ExecuteTrade(g, trade)
}

// NewDecisionManager constructs and links the linear decision processing chain.
func NewDecisionManager() *DecisionManager {
	valNode := &TurnValidationNode{}
	buyNode := &BuyDecisionNode{}
	declineNode := &DeclineBuyDecisionNode{}
	buildNode := &BuildHouseDecisionNode{}
	sellNode := &SellHouseDecisionNode{}
	mortNode := &MortgageDecisionNode{}
	demortNode := &DemortgageDecisionNode{}
	tradeNode := &TradeDecisionNode{}
	bailNode := &PayBailDecisionNode{}

	chain := pipeline.BuildDecisionChain(
		valNode,
		buyNode,
		declineNode,
		buildNode,
		sellNode,
		mortNode,
		demortNode,
		tradeNode,
		bailNode,
	)

	return &DecisionManager{
		chain: chain,
	}
}

func (m *DecisionManager) Execute(g *state.Game) error {
	if m.chain == nil {
		return nil
	}
	return m.chain.Execute(g)
}

func (m *DecisionManager) ExecuteTrade(g *state.Game, trade *TradeOffer) error {
	g.PendingDecision = state.DecisionTrade
	g.PendingTrade = trade
	return m.Execute(g)
}
