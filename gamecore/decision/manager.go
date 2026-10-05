package decision

import (
	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

type DecisionManager struct {
	chain     pipeline.DecisionNode
	tradeNode *TradeDecisionNode
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
		chain:     chain,
		tradeNode: tradeNode,
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
	m.tradeNode.PendingTrade = trade
	return m.Execute(g)
}
