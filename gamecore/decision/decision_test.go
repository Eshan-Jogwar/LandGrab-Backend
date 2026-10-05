package decision

import (
	"errors"
	"testing"

	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/tile"
)

func TestBuyDecision(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	mgr := NewDecisionManager()

	// Landed on Old Harbor Road (pos 1, cost 60)
	g.Phase = state.PhaseAwaitingDecision
	g.PendingDecision = state.DecisionBuy
	g.DecisionTargetTile = 1

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Wallet.Balance() != 1440 {
		t.Errorf("expected 1440 balance, got %d", p1.Wallet.Balance())
	}
	if len(p1.Properties) != 1 {
		t.Fatalf("expected 1 property in player inventory")
	}
	if g.TileOwnerIndex[1] != p1.ID {
		t.Errorf("expected tile 1 owned by p1")
	}
	if g.CurrentPlayer() != p2 {
		t.Errorf("expected turn to advance to p2 after buy decision")
	}
}

func TestDeclineBuyDecision(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})

	mgr := NewDecisionManager()

	g.Phase = state.PhaseAwaitingDecision
	g.PendingDecision = state.DecisionDeclineBuy
	g.DecisionTargetTile = 1

	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Wallet.Balance() != 1500 {
		t.Errorf("expected wallet unchanged, got %d", p1.Wallet.Balance())
	}
	if len(p1.Properties) != 0 {
		t.Errorf("expected no properties bought")
	}
	if g.CurrentPlayer() != p2 {
		t.Errorf("expected turn to advance to p2 after decline decision")
	}
}

func TestBuildHouseMonopolyAndEvenBuildingRule(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})
	mgr := NewDecisionManager()

	// Own only tile 1 (Brown 1)
	p1.AddProperty(tile.TileAt(1))
	g.TileOwnerIndex[1] = p1.ID

	// Try to build on tile 1 without complete monopoly
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err := mgr.Execute(g)
	if !errors.Is(err, ErrMonopolyRequired) {
		t.Fatalf("expected ErrMonopolyRequired, got %v", err)
	}

	// Complete the Brown color group by acquiring tile 3 (Fisherman's Wharf)
	p1.AddProperty(tile.TileAt(3))
	g.TileOwnerIndex[3] = p1.ID

	// Now build 1 house on tile 1 (Cost 50)
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error building 1st house: %v", err)
	}

	ph1, _ := p1.FindProperty(1)
	if ph1.ExtraInfo.Houses != 1 {
		t.Errorf("expected 1 house on tile 1, got %d", ph1.ExtraInfo.Houses)
	}
	if p1.Wallet.Balance() != 1450 {
		t.Errorf("expected 1450 balance after building house, got %d", p1.Wallet.Balance())
	}

	// Try to build a 2nd house on tile 1 before tile 3 has a house (Even building rule violation)
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if !errors.Is(err, ErrEvenBuildingRule) {
		t.Fatalf("expected ErrEvenBuildingRule, got %v", err)
	}

	// Build on tile 3 first
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 3
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error building on tile 3: %v", err)
	}

	ph3, _ := p1.FindProperty(3)
	if ph3.ExtraInfo.Houses != 1 {
		t.Errorf("expected 1 house on tile 3, got %d", ph3.ExtraInfo.Houses)
	}

	// Now tile 1 can build its 2nd house!
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error building 2nd house on tile 1: %v", err)
	}
	if ph1.ExtraInfo.Houses != 2 {
		t.Errorf("expected 2 houses on tile 1, got %d", ph1.ExtraInfo.Houses)
	}
}

func TestBuildHotelAndSell(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})
	mgr := NewDecisionManager()

	ph1 := p1.AddProperty(tile.TileAt(1))
	ph3 := p1.AddProperty(tile.TileAt(3))
	g.TileOwnerIndex[1] = p1.ID
	g.TileOwnerIndex[3] = p1.ID

	// Set both to 4 houses
	ph1.ExtraInfo.Houses = 4
	ph3.ExtraInfo.Houses = 4

	// Upgrade tile 1 to hotel
	g.PendingDecision = state.DecisionBuildHouse
	g.DecisionTargetTile = 1
	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error upgrading to hotel: %v", err)
	}

	if !ph1.ExtraInfo.Hotel || ph1.ExtraInfo.Houses != 0 {
		t.Errorf("expected hotel on tile 1 and houses 0")
	}

	// Selling hotel on tile 1
	g.PendingDecision = state.DecisionSellHouse
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error selling hotel: %v", err)
	}

	if ph1.ExtraInfo.Hotel || ph1.ExtraInfo.Houses != 4 {
		t.Errorf("expected hotel sold and downgraded to 4 houses")
	}
}

func TestMortgageAndDemortgage(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	g := state.NewGame([]*state.Player{p1})
	mgr := NewDecisionManager()

	ph1 := p1.AddProperty(tile.TileAt(1))
	p1.AddProperty(tile.TileAt(3))
	g.TileOwnerIndex[1] = p1.ID
	g.TileOwnerIndex[3] = p1.ID

	// Put 1 house on tile 1
	ph1.ExtraInfo.Houses = 1

	// Attempt to mortgage when houses exist
	g.PendingDecision = state.DecisionMortgage
	g.DecisionTargetTile = 1
	err := mgr.Execute(g)
	if !errors.Is(err, ErrBuildingsInGroup) {
		t.Fatalf("expected ErrBuildingsInGroup, got %v", err)
	}

	// Remove house
	ph1.ExtraInfo.Houses = 0

	// Now mortgage tile 1 (mortgage value = 30)
	g.PendingDecision = state.DecisionMortgage
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error mortgaging: %v", err)
	}

	if !ph1.IsMortgaged {
		t.Errorf("expected tile 1 to be mortgaged")
	}
	if p1.Wallet.Balance() != 1530 {
		t.Errorf("expected balance 1530 after mortgage, got %d", p1.Wallet.Balance())
	}

	// Demortgage tile 1 (mortgage 30 + 10% = 33)
	g.PendingDecision = state.DecisionDemortgage
	g.DecisionTargetTile = 1
	err = mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error demortgaging: %v", err)
	}

	if ph1.IsMortgaged {
		t.Errorf("expected tile 1 to be unmortgaged")
	}
	if p1.Wallet.Balance() != 1497 {
		t.Errorf("expected balance 1497 (1530 - 33), got %d", p1.Wallet.Balance())
	}
}

func TestPayBail(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p1.InJail = true
	g := state.NewGame([]*state.Player{p1})
	mgr := NewDecisionManager()

	g.PendingDecision = state.DecisionPayBail
	err := mgr.Execute(g)
	if err != nil {
		t.Fatalf("unexpected error paying bail: %v", err)
	}

	if p1.InJail {
		t.Errorf("expected InJail to be false after paying bail")
	}
	if p1.Wallet.Balance() != 1450 {
		t.Errorf("expected balance 1450 after $50 bail, got %d", p1.Wallet.Balance())
	}
}

func TestTradeDecision(t *testing.T) {
	p1 := state.NewPlayer(1, 1500)
	p2 := state.NewPlayer(2, 1500)
	g := state.NewGame([]*state.Player{p1, p2})
	mgr := NewDecisionManager()

	t1 := tile.TileAt(1)
	p1.AddProperty(t1)
	g.TileOwnerIndex[1] = p1.ID

	// Trade tile 1 from p1 to p2 for $300
	trade := &TradeOffer{
		FromPlayerID: p1.ID,
		ToPlayerID:   p2.ID,
		TileID:       1,
		Cash:         300,
	}

	err := mgr.ExecuteTrade(g, trade)
	if err != nil {
		t.Fatalf("unexpected error in trade: %v", err)
	}

	if p1.Wallet.Balance() != 1800 {
		t.Errorf("expected p1 balance 1800, got %d", p1.Wallet.Balance())
	}
	if p2.Wallet.Balance() != 1200 {
		t.Errorf("expected p2 balance 1200, got %d", p2.Wallet.Balance())
	}
	if _, ok := p1.FindProperty(1); ok {
		t.Errorf("expected p1 to no longer own tile 1")
	}
	if _, ok := p2.FindProperty(1); !ok {
		t.Errorf("expected p2 to own tile 1")
	}
	if g.TileOwnerIndex[1] != p2.ID {
		t.Errorf("expected TileOwnerIndex to point to p2")
	}
}
