package tile

import (
	"github.com/LandGrab/gamecore/types"
)

type ExtraInfo struct {
	Houses         int
	Hotel          bool
	Dice1          types.Dice
	Dice2          types.Dice
	RailwaysOwned  int
	UtilitiesOwned int
}

type Tile interface {
	ID() types.ID
	Name() string
	Color() Color

	Cost() (types.Money, error)
	Rent(*ExtraInfo) (types.Money, error)

	CostHouse() (types.Money, error)
	CostHotel() (types.Money, error)

	Buildable() bool

	MortgageValue() (types.Money, error)
	DemortgageValue() (types.Money, error)
}
