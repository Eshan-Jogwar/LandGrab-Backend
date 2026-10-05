package pipeline

import (
	"errors"
	"testing"

	"github.com/LandGrab/gamecore/state"
)

type mockMoveNode struct {
	BaseMoveNode
	name          string
	shouldRun     bool
	errToReturn   error
	executedOrder *[]string
}

func (m *mockMoveNode) ShouldExecute(g *state.Game) bool {
	return m.shouldRun
}

func (m *mockMoveNode) Execute(g *state.Game) error {
	if !m.ShouldExecute(g) {
		return m.ExecuteNext(g)
	}
	if m.errToReturn != nil {
		return m.errToReturn
	}
	if m.executedOrder != nil {
		*m.executedOrder = append(*m.executedOrder, m.name)
	}
	return m.ExecuteNext(g)
}

type mockDecisionNode struct {
	BaseDecisionNode
	name          string
	shouldRun     bool
	errToReturn   error
	executedOrder *[]string
}

func (m *mockDecisionNode) ShouldExecute(g *state.Game) bool {
	return m.shouldRun
}

func (m *mockDecisionNode) Execute(g *state.Game) error {
	if !m.ShouldExecute(g) {
		return m.ExecuteNext(g)
	}
	if m.errToReturn != nil {
		return m.errToReturn
	}
	if m.executedOrder != nil {
		*m.executedOrder = append(*m.executedOrder, m.name)
	}
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

func TestDecisionChainExecutionAndSkipping(t *testing.T) {
	var order []string
	d1 := &mockDecisionNode{name: "dec1", shouldRun: true, executedOrder: &order}
	d2 := &mockDecisionNode{name: "dec2", shouldRun: false, executedOrder: &order}
	d3 := &mockDecisionNode{name: "dec3", shouldRun: true, executedOrder: &order}

	head := BuildDecisionChain(d1, d2, d3)
	if head != d1 {
		t.Fatalf("expected head to be d1")
	}

	g := state.NewGame([]*state.Player{state.NewPlayer(1, 1500)})
	err := head.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error executing decision chain: %v", err)
	}

	if len(order) != 2 {
		t.Fatalf("expected 2 nodes executed, got %d", len(order))
	}
	if order[0] != "dec1" || order[1] != "dec3" {
		t.Errorf("expected [dec1, dec3], got %v", order)
	}
}

func TestChainErrorPropagation(t *testing.T) {
	expectedErr := errors.New("pipeline test error")
	n1 := &mockMoveNode{name: "node1", shouldRun: true}
	n2 := &mockMoveNode{name: "node2", shouldRun: true, errToReturn: expectedErr}
	n3 := &mockMoveNode{name: "node3", shouldRun: true}

	head := BuildMoveChain(n1, n2, n3)
	g := state.NewGame([]*state.Player{state.NewPlayer(1, 1500)})
	err := head.Execute(g)

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestSingleNodeChain(t *testing.T) {
	var order []string
	single := &mockMoveNode{name: "only", shouldRun: true, executedOrder: &order}
	head := BuildMoveChain(single)

	g := state.NewGame([]*state.Player{state.NewPlayer(1, 1500)})
	err := head.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 || order[0] != "only" {
		t.Errorf("expected [only], got %v", order)
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
