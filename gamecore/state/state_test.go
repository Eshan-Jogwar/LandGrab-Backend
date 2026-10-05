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
	if _, ok := p.FindProperty(99); ok {
		t.Errorf("expected FindProperty(99) to return false")
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

	// Remove nonexistent property
	if _, ok := p.RemoveProperty(99); ok {
		t.Errorf("expected RemoveProperty(99) to return false")
	}

	// Remove existing property
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

func TestRailwayAndUtilityRentCalculation(t *testing.T) {
	// Railway: Metro Central Station (pos 5)
	rail := tile.TileAt(5)
	phRail := NewPropertyHolder(rail)

	// 1, 2, 3, 4 railways owned -> 25, 50, 100, 200
	expectedRailwayRents := []int{25, 50, 100, 200}
	for i, exp := range expectedRailwayRents {
		rent, err := phRail.CalculateRent(0, 0, i+1, 0)
		if err != nil {
			t.Fatalf("unexpected error for %d railways: %v", i+1, err)
		}
		if int(rent) != exp {
			t.Errorf("expected railway rent %d, got %d", exp, rent)
		}
	}

	// Utility: Power Grid (pos 12)
	util := tile.TileAt(12)
	phUtil := NewPropertyHolder(util)

	// 1 utility, dice 3+4=7 -> 7 * 4 = 28
	rent, err := phUtil.CalculateRent(3, 4, 0, 1)
	if err != nil {
		t.Fatalf("unexpected error for 1 utility: %v", err)
	}
	if rent != 28 {
		t.Errorf("expected utility rent 28, got %d", rent)
	}

	// 2 utilities, dice 5+6=11 -> 11 * 10 = 110
	rent, err = phUtil.CalculateRent(5, 6, 0, 2)
	if err != nil {
		t.Fatalf("unexpected error for 2 utilities: %v", err)
	}
	if rent != 110 {
		t.Errorf("expected utility rent 110, got %d", rent)
	}
}

func TestTotalInColorGroup(t *testing.T) {
	tests := []struct {
		color    tile.Color
		expected int
	}{
		{tile.ColorBrown, 2},
		{tile.ColorDarkBlue, 2},
		{tile.ColorLightBlue, 3},
		{tile.ColorPink, 3},
		{tile.ColorOrange, 3},
		{tile.ColorRed, 3},
		{tile.ColorYellow, 3},
		{tile.ColorGreen, 3},
		{tile.ColorNone, 0},
	}

	for _, tc := range tests {
		got := TotalInColorGroup(tc.color)
		if got != tc.expected {
			t.Errorf("TotalInColorGroup(%v): expected %d, got %d", tc.color, tc.expected, got)
		}
	}
}

func TestGameState(t *testing.T) {
	p1 := NewPlayer(1, 1500)
	p2 := NewPlayer(2, 1500)
	p3 := NewPlayer(3, 1500)
	g := NewGame([]*Player{p1, p2, p3})

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

	// Out of bounds ownership check
	owner, ph = g.GetPropertyOwner(99)
	if owner != nil || ph != nil {
		t.Errorf("expected out of bounds tile to return nil")
	}

	// Assign tile 1 to p1
	t1 := tile.TileAt(1)
	p1Ph := p1.AddProperty(t1)
	g.TileOwnerIndex[1] = p1.ID

	owner, ph = g.GetPropertyOwner(1)
	if owner != p1 || ph != p1Ph {
		t.Errorf("expected GetPropertyOwner(1) to return p1 and p1Ph")
	}

	// Advance turn cycle: p1 -> p2 -> p3 -> p1
	g.AdvanceTurn()
	if g.CurrentPlayer() != p2 {
		t.Errorf("expected p2, got %v", g.CurrentPlayer().ID)
	}
	g.AdvanceTurn()
	if g.CurrentPlayer() != p3 {
		t.Errorf("expected p3, got %v", g.CurrentPlayer().ID)
	}
	g.AdvanceTurn()
	if g.CurrentPlayer() != p1 {
		t.Errorf("expected p1 after full cycle, got %v", g.CurrentPlayer().ID)
	}
}

func TestEmptyGameState(t *testing.T) {
	emptyGame := NewGame(nil)
	if emptyGame.CurrentPlayer() != nil {
		t.Errorf("expected nil CurrentPlayer for empty game")
	}
	if emptyGame.CurrentTile() != nil {
		t.Errorf("expected nil CurrentTile for empty game")
	}
	// Calling AdvanceTurn on empty game should not panic
	emptyGame.AdvanceTurn()
}
