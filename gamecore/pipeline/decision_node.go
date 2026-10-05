package pipeline

import "github.com/LandGrab/gamecore/state"

// DecisionNode represents a link in the decision chain of responsibility.
type DecisionNode interface {
	SetNext(next DecisionNode) MoveDecisionNode
	ShouldExecute(g *state.Game) bool
	Execute(g *state.Game) error
}

// MoveDecisionNode aliases DecisionNode for chaining syntax clarity.
type MoveDecisionNode = DecisionNode

// BaseDecisionNode provides default chaining behavior for decision nodes.
type BaseDecisionNode struct {
	next DecisionNode
}

func (b *BaseDecisionNode) SetNext(next DecisionNode) DecisionNode {
	b.next = next
	return next
}

func (b *BaseDecisionNode) ExecuteNext(g *state.Game) error {
	if b.next != nil {
		return b.next.Execute(g)
	}
	return nil
}
