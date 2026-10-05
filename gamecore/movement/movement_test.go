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

func TestSelfOwnedPropertyLanding(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})

	// Give Old Harbor Road (pos 1) to p1
	p1.AddProperty(tile.TileAt(1))
	g.TileOwnerIndex[1] = p1.ID

	// Roll 0+1 = 1
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 0, 1
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not pay rent to oneself
	if p1.Wallet.Balance() != 1500 {
		t.Errorf("expected balance unchanged at 1500, got %d", p1.Wallet.Balance())
	}
}

func TestMortgagedPropertyLanding(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	ph := p2.AddProperty(tile.TileAt(1))
	ph.IsMortgaged = true
	g.TileOwnerIndex[1] = p2.ID

	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 0, 1
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Rent on mortgaged property is 0
	if p1.Wallet.Balance() != 1500 {
		t.Errorf("expected balance unchanged at 1500 for mortgaged property, got %d", p1.Wallet.Balance())
	}
}

func TestRailwayLandingRent(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Give p2 two railways: pos 5 (Metro Central) and pos 15 (North Shore Line)
	p2.AddProperty(tile.TileAt(5))
	p2.AddProperty(tile.TileAt(15))
	g.TileOwnerIndex[5] = p2.ID
	g.TileOwnerIndex[15] = p2.ID

	// p1 lands on pos 5 (2 railways owned -> rent $50)
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 2, 3 // sum = 5
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Wallet.Balance() != 1450 { // 1500 - 50 = 1450
		t.Errorf("expected p1 wallet 1450, got %d", p1.Wallet.Balance())
	}
	if p2.Wallet.Balance() != 1550 {
		t.Errorf("expected p2 wallet 1550, got %d", p2.Wallet.Balance())
	}
}

func TestUtilityLandingRent(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Give p2 one utility: pos 12 (Power Grid)
	p2.AddProperty(tile.TileAt(12))
	g.TileOwnerIndex[12] = p2.ID

	// p1 lands on pos 12 (roll 6+6=12, 1 utility -> dice sum 12 * 4 = 48)
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 6, 6
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Wallet.Balance() != 1452 { // 1500 - 48 = 1452
		t.Errorf("expected p1 balance 1452, got %d", p1.Wallet.Balance())
	}
	if p2.Wallet.Balance() != 1548 {
		t.Errorf("expected p2 balance 1548, got %d", p2.Wallet.Balance())
	}
}

func TestBankruptcyOnRent(t *testing.T) {
	p1 := state.NewPlayer(1, 10) // Only has $10
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Give p2 Old Harbor Road (rent $20 with 1 house)
	ph := p2.AddProperty(tile.TileAt(1))
	ph.ExtraInfo.Houses = 1 // Rent is 10
	g.TileOwnerIndex[1] = p2.ID

	// Make rent 30 with 2 houses
	ph.ExtraInfo.Houses = 2 // Rent is 30, but p1 only has 10
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 0, 1
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !p1.IsBankrupt {
		t.Errorf("expected p1 to be marked bankrupt")
	}
	if p1.Wallet.Balance() != 0 {
		t.Errorf("expected p1 balance 0 after liquidation, got %d", p1.Wallet.Balance())
	}
	if p2.Wallet.Balance() != 1510 { // Collected all $10 p1 had
		t.Errorf("expected p2 balance 1510, got %d", p2.Wallet.Balance())
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

func TestMovementTaxDeductionAndBankruptcy(t *testing.T) {
	// Case 1: Income Tax ($200 at pos 4)
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})

	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 1, 3 // sum = 4
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

	// Case 2: Wealth Tax ($100 at pos 38)
	p2 := state.NewPlayer(2, 500)
	p2.Position = 36
	g2 := state.NewGame([]*state.Player{p2})

	mgr2 := NewMoveManager(func() (types.Dice, types.Dice) {
		return 1, 1 // sum = 2 -> 36+2=38
	})

	err = mgr2.Execute(g2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p2.Position != 38 {
		t.Errorf("expected pos 38, got %d", p2.Position)
	}
	if p2.Wallet.Balance() != 400 {
		t.Errorf("expected balance 400 after wealth tax, got %d", p2.Wallet.Balance())
	}

	// Case 3: Tax Bankruptcy (cannot pay $200)
	p3 := state.NewPlayer(3, 50)
	g3 := state.NewGame([]*state.Player{p3})

	err = mgr.Execute(g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p3.IsBankrupt {
		t.Errorf("expected p3 to be bankrupt from tax")
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

func TestJailTurnProgressionAndForcedBail(t *testing.T) {
	// Case A: Turn 1 without doubles -> stays in jail
	p1 := state.NewPlayer(1, 1500)
	p1.Position = 10
	p1.InJail = true
	g := state.NewGame([]*state.Player{p1})

	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 1, 2 // non-doubles
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p1.InJail || p1.JailTurnsServed != 1 || p1.Position != 10 {
		t.Errorf("expected turn 1 in jail (turnsServed=1, pos=10)")
	}

	// Case B: Turn 3 with enough funds -> forced $50 bail, released and moves
	p1.JailTurnsServed = 2
	g.Phase = state.PhaseReadyToRoll
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p1.InJail {
		t.Errorf("expected released on 3rd turn after forced bail")
	}
	if p1.Wallet.Balance() != 1450 { // 1500 - 50 = 1450
		t.Errorf("expected $50 deducted for forced bail, got %d", p1.Wallet.Balance())
	}
	if p1.Position != 13 { // 10 + (1+2=3) = 13
		t.Errorf("expected moved to pos 13, got %d", p1.Position)
	}

	// Case C: Turn 3 with insufficient funds -> marked bankrupt
	p2 := state.NewPlayer(2, 20) // Only has $20
	p2.Position = 10
	p2.InJail = true
	p2.JailTurnsServed = 2
	g2 := state.NewGame([]*state.Player{p2})

	err = mgr.Execute(g2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p2.IsBankrupt {
		t.Errorf("expected p2 to be bankrupt from inability to pay forced bail")
	}
}

func TestDoublesGrantsExtraTurn(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	// Roll 3+3=6 (doubles)
	mgr := NewMoveManager(func() (types.Dice, types.Dice) {
		return 3, 3
	})

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Stays on p1 for an extra turn
	if g.CurrentPlayer() != p1 {
		t.Errorf("expected turn to remain on p1 after rolling doubles")
	}
	if g.Phase != state.PhaseAwaitingDecision { // Landed on pos 6 (Crescent Way, unowned)
		t.Errorf("expected PhaseAwaitingDecision on Crescent Way, got %v", g.Phase)
	}
}
