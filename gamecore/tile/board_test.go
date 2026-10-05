package tile

import (
	stdliberrors "errors"
	"testing"

	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

func TestBoardTileCount(t *testing.T) {
	if TileCount() != 40 {
		t.Fatalf("expected TileCount to be 40, got %d", TileCount())
	}

	tiles := AllTiles()
	if len(tiles) != 40 {
		t.Fatalf("expected 40 tiles in AllTiles, got %d", len(tiles))
	}

	// Verify all tiles are non-nil
	for i, tileRef := range tiles {
		if tileRef == nil {
			t.Errorf("tile at index %d is nil", i)
		}
	}
}

func TestBoardCircularWrapping(t *testing.T) {
	tests := []struct {
		pos          int
		expectedName string
		expectedID   types.ID
	}{
		{0, "Start", 1},
		{40, "Start", 1},
		{80, "Start", 1},
		{-40, "Start", 1},
		{1, "Old Harbor Road", 2},
		{41, "Old Harbor Road", 2},
		{39, "Celestial Heights", 40},
		{79, "Celestial Heights", 40},
		{-1, "Celestial Heights", 40},
		{-41, "Celestial Heights", 40},
		{5, "Metro Central Station", 6},
		{45, "Metro Central Station", 6},
	}

	for _, tc := range tests {
		tRef := TileAt(tc.pos)
		if tRef == nil {
			t.Fatalf("TileAt(%d) returned nil", tc.pos)
		}
		if tRef.Name() != tc.expectedName {
			t.Errorf("TileAt(%d): expected name %q, got %q", tc.pos, tc.expectedName, tRef.Name())
		}
		if tRef.ID() != tc.expectedID {
			t.Errorf("TileAt(%d): expected ID %d, got %d", tc.pos, tc.expectedID, tRef.ID())
		}
	}
}

func TestBoardFlyweightIdentity(t *testing.T) {
	// Verify that multiple calls return the exact same flyweight instance
	t1 := TileAt(1)
	t2 := TileAt(41) // 1 + 40
	t3 := TileAt(-39)

	if t1 != t2 {
		t.Errorf("expected t1 and t2 to have identical pointer address, got %p and %p", t1, t2)
	}
	if t1 != t3 {
		t.Errorf("expected t1 and t3 to have identical pointer address, got %p and %p", t1, t3)
	}

	byID, err := TileByID(2)
	if err != nil {
		t.Fatalf("unexpected error looking up ID 2: %v", err)
	}
	if byID != t1 {
		t.Errorf("expected TileByID(2) to return identical instance as TileAt(1)")
	}
}

func TestBoardTileByID(t *testing.T) {
	for id := types.ID(1); id <= 40; id++ {
		tileRef, err := TileByID(id)
		if err != nil {
			t.Fatalf("TileByID(%d) returned unexpected error: %v", id, err)
		}
		if tileRef.ID() != id {
			t.Errorf("expected tile ID %d, got %d", id, tileRef.ID())
		}
	}

	// Invalid IDs
	invalidIDs := []types.ID{0, 41, 100}
	for _, id := range invalidIDs {
		_, err := TileByID(id)
		if err == nil {
			t.Errorf("expected error for invalid ID %d, got nil", id)
		}
		if !stdliberrors.Is(err, errors.ErrTileNotFound) {
			t.Errorf("expected ErrTileNotFound for ID %d, got %v", id, err)
		}
	}
}

func TestBoardColorGroups(t *testing.T) {
	// Test that street properties have their colors
	brown1 := TileAt(1)
	if brown1.Color() != ColorBrown {
		t.Errorf("expected Old Harbor Road to be ColorBrown, got %v", brown1.Color())
	}

	darkBlue2 := TileAt(39)
	if darkBlue2.Color() != ColorDarkBlue {
		t.Errorf("expected Celestial Heights to be ColorDarkBlue, got %v", darkBlue2.Color())
	}

	// Test that non-properties have ColorNone
	start := TileAt(0)
	if start.Color() != ColorNone {
		t.Errorf("expected Start to have ColorNone, got %v", start.Color())
	}

	rail := TileAt(5)
	if rail.Color() != ColorNone {
		t.Errorf("expected Railway to have ColorNone, got %v", rail.Color())
	}

	utility := TileAt(12)
	if utility.Color() != ColorNone {
		t.Errorf("expected Utility to have ColorNone, got %v", utility.Color())
	}
}

func TestBoardImmutability(t *testing.T) {
	// Modifying the slice returned by AllTiles does not affect the board
	tiles := AllTiles()
	tiles[0] = nil

	if TileAt(0) == nil {
		t.Errorf("expected defaultBoard[0] to remain intact after mutating copy")
	}
}
