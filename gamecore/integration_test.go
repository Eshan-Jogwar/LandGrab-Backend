package gamecore_test

import (
	"testing"

	"github.com/LandGrab/gamecore/decision"
	"github.com/LandGrab/gamecore/movement"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/types"
)

func TestFullGameTurnLoop(t *testing.T) {
	// Initialize Alice (ID 1) and Bob (ID 2)
	alice := state.NewPlayer(1, 1500)
	bob := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{alice, bob})

	decMgr := decision.NewDecisionManager()

	// MoveManager with controllable dice
	var nextRoll [2]types.Dice
	moveMgr := movement.NewMoveManager(func() (types.Dice, types.Dice) {
		return nextRoll[0], nextRoll[1]
	})

	// ---------------------------------------------------------
	// Turn 1: Alice moves to tile 1 (Old Harbor Road, unowned)
	// ---------------------------------------------------------
	nextRoll = [2]types.Dice{0, 1}
	err := moveMgr.Execute(g)
	if err != nil {
		t.Fatalf("turn 1 move error: %v", err)
	}

	if g.Phase != state.PhaseAwaitingDecision || g.PendingDecision != state.DecisionBuy {
		t.Fatalf("expected PhaseAwaitingDecision & DecisionBuy, got %v / %v", g.Phase, g.PendingDecision)
	}

	// Alice buys Old Harbor Road (cost: 60)
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("turn 1 buy error: %v", err)
	}

	if alice.Wallet.Balance() != 1440 {
		t.Errorf("expected Alice balance 1440, got %d", alice.Wallet.Balance())
	}
	if g.CurrentPlayer() != bob {
		t.Errorf("expected turn to advance to Bob")
	}

	// ---------------------------------------------------------
	// Turn 2: Bob moves to tile 1 (Alice's property) -> pays $2 rent
	// ---------------------------------------------------------
	nextRoll = [2]types.Dice{0, 1}
	err = moveMgr.Execute(g)
	if err != nil {
		t.Fatalf("turn 2 move error: %v", err)
	}

	// Base rent is $2
	if bob.Wallet.Balance() != 1498 {
		t.Errorf("expected Bob balance 1498, got %d", bob.Wallet.Balance())
	}
	if alice.Wallet.Balance() != 1442 {
		t.Errorf("expected Alice balance 1442, got %d", alice.Wallet.Balance())
	}
	if g.CurrentPlayer() != alice {
		t.Errorf("expected turn to advance to Alice")
	}

	// ---------------------------------------------------------
	// Turn 3: Alice moves from tile 1 to tile 3 (Fisherman's Wharf, unowned)
	// ---------------------------------------------------------
	nextRoll = [2]types.Dice{1, 1} // roll 2 -> 1 + 2 = 3 (doubles!)
	err = moveMgr.Execute(g)
	if err != nil {
		t.Fatalf("turn 3 move error: %v", err)
	}

	// Alice buys tile 3 (cost: 60) -> completing Brown monopoly!
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("turn 3 buy error: %v", err)
	}

	if alice.Wallet.Balance() != 1382 { // 1442 - 60 = 1382
		t.Errorf("expected Alice balance 1382, got %d", alice.Wallet.Balance())
	}

	// Because Alice rolled doubles, it is still Alice's turn!
	if g.CurrentPlayer() != alice {
		t.Errorf("expected Alice to have another roll due to doubles")
	}

	// ---------------------------------------------------------
	// Building Phase: Alice develops houses on Brown monopoly
	// ---------------------------------------------------------
	// 1. Build 1 house on tile 1 (cost 50)
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("build house 1 on tile 1 error: %v", err)
	}

	// 2. Build 1 house on tile 3 (cost 50)
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 3
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("build house 1 on tile 3 error: %v", err)
	}

	// 3. Build 2nd house on tile 1 (cost 50)
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("build house 2 on tile 1 error: %v", err)
	}

	// 1382 - 150 = 1232
	if alice.Wallet.Balance() != 1232 {
		t.Errorf("expected Alice balance 1232 after 3 houses, got %d", alice.Wallet.Balance())
	}

	// Alice finishes her extra turn with a non-doubles roll
	nextRoll = [2]types.Dice{2, 3} // sum = 5 -> pos 3 + 5 = 8 (Beacon Street, unowned)
	err = moveMgr.Execute(g)
	if err != nil {
		t.Fatalf("alice extra roll error: %v", err)
	}
	// Alice declines buying Beacon Street -> turn advances to Bob
	g.PendingDecision = state.DecisionDeclineBuy
	err = decMgr.Execute(g)
	if err != nil {
		t.Fatalf("alice decline error: %v", err)
	}
	if g.CurrentPlayer() != bob {
		t.Errorf("expected turn to advance to Bob")
	}

	// ---------------------------------------------------------
	// Turn 4: Bob wraps around the board and lands on Alice's 2-house property
	// ---------------------------------------------------------
	// Bob is at pos 1. Rolls 40 (wraps 1 + 40 = 41 -> pos 1, passing Go!)
	nextRoll = [2]types.Dice{20, 20}
	err = moveMgr.Execute(g)
	if err != nil {
		t.Fatalf("bob wrap error: %v", err)
	}

	// Bob passed Go (+200: 1498 + 200 = 1698)
	// Bob landed on tile 1 (2 houses rent = $30: 1698 - 30 = 1668)
	// Alice received $30 rent (1232 + 30 = 1262)
	if bob.Wallet.Balance() != 1668 {
		t.Errorf("expected Bob balance 1668, got %d", bob.Wallet.Balance())
	}
	if alice.Wallet.Balance() != 1262 {
		t.Errorf("expected Alice balance 1262, got %d", alice.Wallet.Balance())
	}
}
