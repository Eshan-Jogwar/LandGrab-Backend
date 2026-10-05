package tile

import (
	"errors"
	"testing"

	gameerrors "github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

func TestPropertyTile(t *testing.T) {
	p := NewProperty(PropertyConfig{
		ID:            1,
		Name:          "Boardwalk",
		Color:         ColorDarkBlue,
		Cost:          400,
		Rent:          50,
		Rent1House:    200,
		Rent2House:    600,
		Rent3House:    1400,
		Rent4House:    1700,
		RentHotel:     2000,
		CostHouse:     200,
		CostHotel:     200,
		MortgageValue: 200,
	})

	if p.Name() != "Boardwalk" {
		t.Errorf("expected name 'Boardwalk', got %s", p.Name())
	}

	if prop, ok := p.(interface{ Color() Color }); !ok || prop.Color() != ColorDarkBlue {
		t.Errorf("expected color ColorDarkBlue, got %v", prop.Color())
	}

	cost, err := p.Cost()
	if err != nil || cost != 400 {
		t.Errorf("expected cost 400, nil; got %v, %v", cost, err)
	}

	if !p.Buildable() {
		t.Errorf("expected property to be buildable")
	}

	houseCost, err := p.CostHouse()
	if err != nil || houseCost != 200 {
		t.Errorf("expected house cost 200, nil; got %v, %v", houseCost, err)
	}

	hotelCost, err := p.CostHotel()
	if err != nil || hotelCost != 200 {
		t.Errorf("expected hotel cost 200, nil; got %v, %v", hotelCost, err)
	}

	mortgage, err := p.MortgageValue()
	if err != nil || mortgage != 200 {
		t.Errorf("expected mortgage 200, nil; got %v, %v", mortgage, err)
	}

	demortgage, err := p.DemortgageValue()
	if err != nil || demortgage != 220 {
		t.Errorf("expected demortgage 220, nil; got %v, %v", demortgage, err)
	}

	tests := []struct {
		name         string
		info         *ExtraInfo
		expectedRent types.Money
		expectErr    bool
	}{
		{
			name:      "nil extra info",
			info:      nil,
			expectErr: true,
		},
		{
			name:         "base rent (0 houses)",
			info:         &ExtraInfo{Houses: 0},
			expectedRent: 50,
		},
		{
			name:         "1 house",
			info:         &ExtraInfo{Houses: 1},
			expectedRent: 200,
		},
		{
			name:         "2 houses",
			info:         &ExtraInfo{Houses: 2},
			expectedRent: 600,
		},
		{
			name:         "3 houses",
			info:         &ExtraInfo{Houses: 3},
			expectedRent: 1400,
		},
		{
			name:         "4 houses",
			info:         &ExtraInfo{Houses: 4},
			expectedRent: 1700,
		},
		{
			name:         "hotel",
			info:         &ExtraInfo{Hotel: true},
			expectedRent: 2000,
		},
		{
			name:      "negative houses",
			info:      &ExtraInfo{Houses: -1},
			expectErr: true,
		},
		{
			name:      "too many houses",
			info:      &ExtraInfo{Houses: 5},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rent, err := p.Rent(tc.info)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, gameerrors.ErrInvalidExtraInfo) {
					t.Errorf("expected ErrInvalidExtraInfo, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rent != tc.expectedRent {
					t.Errorf("expected rent %v, got %v", tc.expectedRent, rent)
				}
			}
		})
	}
}

func TestRailwayTile(t *testing.T) {
	rw := NewRailway(2, "Reading Railroad", 200, 25, 100)

	if rw.Name() != "Reading Railroad" {
		t.Errorf("expected name 'Reading Railroad', got %s", rw.Name())
	}

	cost, err := rw.Cost()
	if err != nil || cost != 200 {
		t.Errorf("expected cost 200, nil; got %v, %v", cost, err)
	}

	if rw.Buildable() {
		t.Errorf("expected railway not to be buildable")
	}

	if _, err := rw.CostHouse(); !errors.Is(err, gameerrors.ErrNotBuildable) {
		t.Errorf("expected CostHouse to return ErrNotBuildable, got %v", err)
	}

	if _, err := rw.CostHotel(); !errors.Is(err, gameerrors.ErrNotBuildable) {
		t.Errorf("expected CostHotel to return ErrNotBuildable, got %v", err)
	}

	mortgage, err := rw.MortgageValue()
	if err != nil || mortgage != 100 {
		t.Errorf("expected mortgage 100, nil; got %v, %v", mortgage, err)
	}

	demortgage, err := rw.DemortgageValue()
	if err != nil || demortgage != 110 {
		t.Errorf("expected demortgage 110, nil; got %v, %v", demortgage, err)
	}

	rentTests := []struct {
		name         string
		info         *ExtraInfo
		expectedRent types.Money
		expectErr    bool
	}{
		{
			name:      "nil extra info",
			info:      nil,
			expectErr: true,
		},
		{
			name:      "zero railways owned",
			info:      &ExtraInfo{RailwaysOwned: 0},
			expectErr: true,
		},
		{
			name:         "1 railway owned",
			info:         &ExtraInfo{RailwaysOwned: 1},
			expectedRent: 25,
		},
		{
			name:         "2 railways owned",
			info:         &ExtraInfo{RailwaysOwned: 2},
			expectedRent: 50,
		},
		{
			name:         "3 railways owned",
			info:         &ExtraInfo{RailwaysOwned: 3},
			expectedRent: 100,
		},
		{
			name:         "4 railways owned",
			info:         &ExtraInfo{RailwaysOwned: 4},
			expectedRent: 200,
		},
	}

	for _, tc := range rentTests {
		t.Run(tc.name, func(t *testing.T) {
			rent, err := rw.Rent(tc.info)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, gameerrors.ErrInvalidExtraInfo) {
					t.Errorf("expected ErrInvalidExtraInfo, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rent != tc.expectedRent {
					t.Errorf("expected rent %v, got %v", tc.expectedRent, rent)
				}
			}
		})
	}
}

func TestUtilityTile(t *testing.T) {
	u := NewUtility(3, "Electric Company", 150, 75)

	if u.Name() != "Electric Company" {
		t.Errorf("expected name 'Electric Company', got %s", u.Name())
	}

	cost, err := u.Cost()
	if err != nil || cost != 150 {
		t.Errorf("expected cost 150, nil; got %v, %v", cost, err)
	}

	if u.Buildable() {
		t.Errorf("expected utility not to be buildable")
	}

	if _, err := u.CostHouse(); !errors.Is(err, gameerrors.ErrNotBuildable) {
		t.Errorf("expected CostHouse to return ErrNotBuildable, got %v", err)
	}

	if _, err := u.CostHotel(); !errors.Is(err, gameerrors.ErrNotBuildable) {
		t.Errorf("expected CostHotel to return ErrNotBuildable, got %v", err)
	}

	mortgage, err := u.MortgageValue()
	if err != nil || mortgage != 75 {
		t.Errorf("expected mortgage 75, nil; got %v, %v", mortgage, err)
	}

	demortgage, err := u.DemortgageValue()
	if err != nil || demortgage != 82 { // 75 + 7 = 82
		t.Errorf("expected demortgage 82, nil; got %v, %v", demortgage, err)
	}

	rentTests := []struct {
		name         string
		info         *ExtraInfo
		expectedRent types.Money
		expectErr    bool
	}{
		{
			name:      "nil extra info",
			info:      nil,
			expectErr: true,
		},
		{
			name:      "zero utilities owned",
			info:      &ExtraInfo{UtilitiesOwned: 0, Dice1: 3, Dice2: 4},
			expectErr: true,
		},
		{
			name:         "1 utility owned (dice 3+4=7, 4x)",
			info:         &ExtraInfo{UtilitiesOwned: 1, Dice1: 3, Dice2: 4},
			expectedRent: 28,
		},
		{
			name:         "2 utilities owned (dice 5+6=11, 10x)",
			info:         &ExtraInfo{UtilitiesOwned: 2, Dice1: 5, Dice2: 6},
			expectedRent: 110,
		},
	}

	for _, tc := range rentTests {
		t.Run(tc.name, func(t *testing.T) {
			rent, err := u.Rent(tc.info)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, gameerrors.ErrInvalidExtraInfo) {
					t.Errorf("expected ErrInvalidExtraInfo, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rent != tc.expectedRent {
					t.Errorf("expected rent %v, got %v", tc.expectedRent, rent)
				}
			}
		})
	}
}

func TestNormalTiles(t *testing.T) {
	normalTileNames := []string{
		"Go",
		"Income Tax",
		"Chance",
		"Community Chest",
		"Jail",
		"Free Parking",
		"Go To Jail",
	}

	for id, name := range normalTileNames {
		t.Run(name, func(t *testing.T) {
			tile := NewNormal(types.ID(id+100), name)

			if tile.Name() != name {
				t.Errorf("expected name %q, got %q", name, tile.Name())
			}

			if tile.Buildable() {
				t.Errorf("expected normal tile to not be buildable")
			}

			if _, err := tile.Cost(); !errors.Is(err, gameerrors.ErrNotPurchasable) {
				t.Errorf("expected ErrNotPurchasable, got %v", err)
			}

			if _, err := tile.Rent(&ExtraInfo{}); !errors.Is(err, gameerrors.ErrNotRentable) {
				t.Errorf("expected ErrNotRentable, got %v", err)
			}

			if _, err := tile.CostHouse(); !errors.Is(err, gameerrors.ErrNotBuildable) {
				t.Errorf("expected ErrNotBuildable, got %v", err)
			}

			if _, err := tile.CostHotel(); !errors.Is(err, gameerrors.ErrNotBuildable) {
				t.Errorf("expected CostHotel to return ErrNotBuildable, got %v", err)
			}

			if _, err := tile.MortgageValue(); !errors.Is(err, gameerrors.ErrNoMortgage) {
				t.Errorf("expected ErrNoMortgage, got %v", err)
			}

			if _, err := tile.DemortgageValue(); !errors.Is(err, gameerrors.ErrNoMortgage) {
				t.Errorf("expected ErrNoMortgage, got %v", err)
			}
		})
	}
}
