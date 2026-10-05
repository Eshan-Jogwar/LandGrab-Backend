package pipeline

import (
	"testing"

	"github.com/LandGrab/gamecore/state"
)

type mockMoveNode struct {
	BaseMoveNode
	name          string
	shouldRun     bool
	executedOrder *[]string
}

func (m *mockMoveNode) ShouldExecute(g *state.Game) bool {
	return m.shouldRun
}

func (m *mockMoveNode) Execute(g *state.Game) error {
	if !m.ShouldExecute(g) {
		return m.ExecuteNext(g)
	}
	*m.executedOrder = append(*m.executedOrder, m.name)
	return m.ExecuteNext(g)
}

func TestMoveChainExecutionAndSkipping(t *testing.T) {
	var order []string
	n1 := &mockMoveNode{name: "node1", shouldRun: true, executedOrder: &order}
	n2 := &mockMoveNode{name: "node2", shouldRun: false, executedOrder: &order}
	n3 := &mockMoveNode{name: "node3", shouldRun: true, executedOrder: &order}

	head := BuildMoveChain(n1, n2, n3)
	if head != n1 {
		t.Fatalf("expected head to be n1")
	}

	g := state.NewGame([]*state.Player{state.NewPlayer(1, 1500)})
	err := head.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error executing chain: %v", err)
	}

	if len(order) != 2 {
		t.Fatalf("expected 2 nodes executed, got %d (%v)", len(order), order)
	}
	if order[0] != "node1" || order[1] != "node3" {
		t.Errorf("expected [node1, node3], got %v", order)
	}
}

func TestBuildChainEmpty(t *testing.T) {
	if BuildMoveChain() != nil {
		t.Errorf("expected empty MoveNode chain to return nil")
	}
	if BuildDecisionChain() != nil {
		t.Errorf("expected empty DecisionNode chain to return nil")
	}
}
