package pipeline

import "github.com/LandGrab/gamecore/state"

// MoveNode represents a link in the movement chain of responsibility.
type MoveNode interface {
	SetNext(next MoveNode) MoveNode
	ShouldExecute(g *state.Game) bool
	Execute(g *state.Game) error
}

// BaseMoveNode provides default chaining behavior for movement nodes.
type BaseMoveNode struct {
	next MoveNode
}

func (b *BaseMoveNode) SetNext(next MoveNode) MoveNode {
	b.next = next
	return next
}

func (b *BaseMoveNode) ExecuteNext(g *state.Game) error {
	if b.next != nil {
		return b.next.Execute(g)
	}
	return nil
}
