package state

import (
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

// PropertyHolder bundles the static flyweight tile with runtime ownership state.
type PropertyHolder struct {
	Tile        tile.Tile
	IsMortgaged bool
	ExtraInfo   tile.ExtraInfo
}

// NewPropertyHolder creates a new PropertyHolder for an un-mortgaged, freshly acquired tile.
func NewPropertyHolder(t tile.Tile) *PropertyHolder {
	return &PropertyHolder{
		Tile:        t,
		IsMortgaged: false,
		ExtraInfo:   tile.ExtraInfo{},
	}
}

// CalculateRent computes the rent due by passing dynamic ExtraInfo to the flyweight tile.
// If the property is mortgaged, 0 rent is returned.
func (ph *PropertyHolder) CalculateRent(dice1, dice2 types.Dice, railwaysOwned, utilitiesOwned int) (types.Money, error) {
	if ph.IsMortgaged {
		return 0, nil
	}
	info := ph.ExtraInfo
	info.Dice1 = dice1
	info.Dice2 = dice2
	info.RailwaysOwned = railwaysOwned
	info.UtilitiesOwned = utilitiesOwned
	return ph.Tile.Rent(&info)
}
