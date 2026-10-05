package movement

import (
	"testing"

	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

func TestMovementToUnownedProperty(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Inject fixed roll of 0+1 = 1 (Old Harbor Road)
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 0, 1
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Position != 1 {
		t.Errorf("expected p1 position 1, got %d", p1.Position)
	}
	if g.Phase != state.PhaseAwaitingDecision {
		t.Errorf("expected PhaseAwaitingDecision, got %v", g.Phase)
	}
	if g.PendingDecision != state.DecisionBuy {
		t.Errorf("expected DecisionBuy, got %v", g.PendingDecision)
	}
	if g.DecisionTargetTile != 1 {
		t.Errorf("expected target tile 1, got %d", g.DecisionTargetTile)
	}
}

func TestMovementRentPayment(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Give Old Harbor Road (pos 1, base rent 2) to p2
	t1 := tile.TileAt(1)
	p2.AddProperty(t1)
	g.TileOwnerIndex[1] = p2.ID

	// Roll 1 for p1
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 0, 1
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Wallet.Balance() != 1498 {
		t.Errorf("expected p1 wallet 1498, got %d", p1.Wallet.Balance())
	}
	if p2.Wallet.Balance() != 1502 {
		t.Errorf("expected p2 wallet 1502, got %d", p2.Wallet.Balance())
	}
}

func TestMovementPassingGo(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p1.Position = 38
	g := state.NewGame([]*state.Player{p1})

	// Roll 1+3=4 (from 38 + 4 = 42 -> wraps to 2, passing Go!)
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 1, 3
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Position != 2 {
		t.Errorf("expected position 2, got %d", p1.Position)
	}
	// 1500 + 200 = 1700
	if p1.Wallet.Balance() != 1700 {
		t.Errorf("expected balance 1700 after passing Go, got %d", p1.Wallet.Balance())
	}
}

func TestMovementTaxDeduction(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})

	// Roll 2+2=4 -> Income Tax (pos 4, costs $200)
	// Notice 2+2 is doubles, but tax should still deduct $200
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 2, 2
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Position != 4 {
		t.Errorf("expected pos 4, got %d", p1.Position)
	}
	if p1.Wallet.Balance() != 1300 {
		t.Errorf("expected balance 1300, got %d", p1.Wallet.Balance())
	}
}

func TestGoToJailTile(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p1.Position = 25
	g := state.NewGame([]*state.Player{p1})

	// Roll 2+3=5 (from 25 to 30: "Go To Prison")
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 2, 3
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Position != 10 {
		t.Errorf("expected teleported to jail (pos 10), got %d", p1.Position)
	}
	if !p1.InJail {
		t.Errorf("expected InJail to be true")
	}
}

func TestThreeDoublesSendToJail(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})
	g.ConsecutiveDoubles = 2 // Already rolled 2 doubles

	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 4, 4 // 3rd consecutive doubles!
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !p1.InJail || p1.Position != 10 {
		t.Errorf("expected 3rd doubles to send to jail (pos 10, inJail true)")
	}
}

func TestJailReleaseOnDoubles(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p1.Position = 10
	p1.InJail = true
	g := state.NewGame([]*state.Player{p1})

	// Roll doubles 2+2=4 while in jail
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 2, 2
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.InJail {
		t.Errorf("expected player to be released from jail on doubles")
	}
	// Moved 10 + 4 = 14 (Lakeview Drive)
	if p1.Position != 14 {
		t.Errorf("expected position 14 after exiting jail, got %d", p1.Position)
	}
}
