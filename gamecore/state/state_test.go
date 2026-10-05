package state

import (
	"testing"

	"github.com/LandGrab/gamecore/tile"
)

func TestWallet(t *testing.T) {
	w := Wallet{}
	w.Add(1500)
	if w.Balance() != 1500 {
		t.Fatalf("expected 1500, got %d", w.Balance())
	}
	if !w.Has(500) {
		t.Errorf("expected Has(500) to be true")
	}
	if w.Has(2000) {
		t.Errorf("expected Has(2000) to be false")
	}
	if !w.Deduct(500) {
		t.Errorf("expected Deduct(500) to succeed")
	}
	if w.Balance() != 1000 {
		t.Errorf("expected 1000 balance, got %d", w.Balance())
	}
	if w.Deduct(2000) {
		t.Errorf("expected Deduct(2000) to fail")
	}
	if w.Balance() != 1000 {
		t.Errorf("expected balance unchanged at 1000, got %d", w.Balance())
	}
}

func TestPlayerAndProperties(t *testing.T) {
	p := NewPlayer(1, 1500)
	if p.ID != 1 {
		t.Errorf("expected ID 1, got %d", p.ID)
	}
	if p.Wallet.Balance() != 1500 {
		t.Errorf("expected 1500 balance, got %d", p.Wallet.Balance())
	}

	// Add property
	brown1 := tile.TileAt(1) // Old Harbor Road
	ph := p.AddProperty(brown1)
	if ph == nil {
		t.Fatalf("AddProperty returned nil")
	}
	if ph.Tile != brown1 {
		t.Errorf("expected ph.Tile to be brown1")
	}

	// Find property
	found, ok := p.FindProperty(1)
	if !ok || found != ph {
		t.Errorf("expected FindProperty(1) to find ph")
	}

	// Count color group
	if p.CountColorGroup(tile.ColorBrown) != 1 {
		t.Errorf("expected 1 brown property, got %d", p.CountColorGroup(tile.ColorBrown))
	}

	// Add railway (pos 5: Metro Central Station)
	rail1 := tile.TileAt(5)
	p.AddProperty(rail1)
	if p.CountRailways() != 1 {
		t.Errorf("expected 1 railway, got %d", p.CountRailways())
	}

	// Add utility (pos 12: Power Grid)
	util1 := tile.TileAt(12)
	p.AddProperty(util1)
	if p.CountUtilities() != 1 {
		t.Errorf("expected 1 utility, got %d", p.CountUtilities())
	}

	// Remove property
	removed, ok := p.RemoveProperty(1)
	if !ok || removed != ph {
		t.Errorf("expected RemoveProperty(1) to succeed")
	}
	if _, ok := p.FindProperty(1); ok {
		t.Errorf("expected property 1 to be removed")
	}
}

func TestPropertyHolderRentCalculation(t *testing.T) {
	brown1 := tile.TileAt(1) // Old Harbor Road: Base rent 2
	ph := NewPropertyHolder(brown1)

	// Unmortgaged base rent
	rent, err := ph.CalculateRent(0, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rent != 2 {
		t.Errorf("expected base rent 2, got %d", rent)
	}

	// Add houses
	ph.ExtraInfo.Houses = 1
	rent, err = ph.CalculateRent(0, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rent != 10 { // Rent1House is 10
		t.Errorf("expected 1-house rent 10, got %d", rent)
	}

	// Mortgaged rent must be 0
	ph.IsMortgaged = true
	rent, err = ph.CalculateRent(0, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rent != 0 {
		t.Errorf("expected mortgaged rent 0, got %d", rent)
	}
}

func TestGameState(t *testing.T) {
	p1 := NewPlayer(1, 1500)
	p2 := NewPlayer(2, 1500)
	g := NewGame([]*Player{p1, p2})

	if g.CurrentPlayer() != p1 {
		t.Errorf("expected CurrentPlayer to be p1")
	}
	if g.CurrentTile().ID() != 0 {
		t.Errorf("expected initial tile to be Start (ID 0), got %d", g.CurrentTile().ID())
	}

	// Ownership check
	owner, ph := g.GetPropertyOwner(1)
	if owner != nil || ph != nil {
		t.Errorf("expected tile 1 to be unowned initially")
	}

	// Assign tile 1 to p1
	t1 := tile.TileAt(1)
	p1Ph := p1.AddProperty(t1)
	g.TileOwnerIndex[1] = p1.ID

	owner, ph = g.GetPropertyOwner(1)
	if owner != p1 || ph != p1Ph {
		t.Errorf("expected GetPropertyOwner(1) to return p1 and p1Ph")
	}

	// Advance turn
	g.AdvanceTurn()
	if g.CurrentPlayer() != p2 {
		t.Errorf("expected CurrentPlayer after turn advance to be p2")
	}
}
