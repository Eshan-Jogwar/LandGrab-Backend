package state

import (
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

type Wallet struct {
	money types.Money
}

func (w *Wallet) Add(amount types.Money) {
	w.money += amount
}

func (w *Wallet) Deduct(amount types.Money) bool {
	if w.money < amount {
		return false
	}
	w.money -= amount
	return true
}

func (w *Wallet) Balance() types.Money {
	return w.money
}

func (w *Wallet) Has(amount types.Money) bool {
	return w.money >= amount
}

type Player struct {
	ID              types.ID
	Wallet          Wallet
	Properties      []*PropertyHolder
	Position        int
	InJail          bool
	JailTurnsServed int
	IsBankrupt      bool
}

func NewPlayer(id types.ID, startingMoney types.Money) *Player {
	return &Player{
		ID:         id,
		Wallet:     Wallet{money: startingMoney},
		Properties: make([]*PropertyHolder, 0),
		Position:   0,
	}
}

func (p *Player) FindProperty(tileID types.ID) (*PropertyHolder, bool) {
	for _, ph := range p.Properties {
		if ph.Tile.ID() == tileID {
			return ph, true
		}
	}
	return nil, false
}

func (p *Player) AddProperty(t tile.Tile) *PropertyHolder {
	ph := NewPropertyHolder(t)
	p.Properties = append(p.Properties, ph)
	return ph
}

func (p *Player) RemoveProperty(tileID types.ID) (*PropertyHolder, bool) {
	for i, ph := range p.Properties {
		if ph.Tile.ID() == tileID {
			p.Properties = append(p.Properties[:i], p.Properties[i+1:]...)
			return ph, true
		}
	}
	return nil, false
}

func (p *Player) CountColorGroup(color tile.Color) int {
	count := 0
	for _, ph := range p.Properties {
		if ph.Tile.Color() == color {
			count++
		}
	}
	return count
}

func (p *Player) CountRailways() int {
	count := 0
	for _, ph := range p.Properties {
		if IsRailway(ph.Tile) {
			count++
		}
	}
	return count
}

func (p *Player) CountUtilities() int {
	count := 0
	for _, ph := range p.Properties {
		if IsUtility(ph.Tile) {
			count++
		}
	}
	return count
}

// Tile category helpers
func IsStreetProperty(t tile.Tile) bool {
	return t.Buildable()
}

func IsRailway(t tile.Tile) bool {
	_, err := t.Cost()
	return err == nil && !t.Buildable() && !IsUtility(t)
}

func IsUtility(t tile.Tile) bool {
	return t.Name() == "Power Grid" || t.Name() == "Water Works"
}

func IsPurchasable(t tile.Tile) bool {
	_, err := t.Cost()
	return err == nil
}

// TotalInColorGroup returns the number of properties required to complete a monopoly in that color.
func TotalInColorGroup(color tile.Color) int {
	switch color {
	case tile.ColorBrown, tile.ColorDarkBlue:
		return 2
	case tile.ColorLightBlue, tile.ColorPink, tile.ColorOrange, tile.ColorRed, tile.ColorYellow, tile.ColorGreen:
		return 3
	default:
		return 0
	}
}
