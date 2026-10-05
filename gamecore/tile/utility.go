package tile

import (
	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

type utility struct {
	id            types.ID
	name          string
	cost          types.Money
	mortgageValue types.Money
}

func NewUtility(id types.ID, name string, cost types.Money, mortgageValue types.Money) Tile {
	return &utility{
		id:            id,
		name:          name,
		cost:          cost,
		mortgageValue: mortgageValue,
	}
}

func (u *utility) ID() types.ID {
	return u.id
}

func (u *utility) Color() Color {
	return ColorNone
}

func (u *utility) Name() string {
	return u.name
}

func (u *utility) Cost() (types.Money, error) {
	return u.cost, nil
}

func (u *utility) Rent(info *ExtraInfo) (types.Money, error) {
	if info == nil || info.UtilitiesOwned <= 0 {
		return 0, errors.ErrInvalidExtraInfo
	}
	diceSum := types.Money(info.Dice1 + info.Dice2)
	if info.UtilitiesOwned == 1 {
		return diceSum * 4, nil
	}
	return diceSum * 10, nil
}

func (u *utility) CostHouse() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (u *utility) CostHotel() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (u *utility) Buildable() bool {
	return false
}

func (u *utility) MortgageValue() (types.Money, error) {
	return u.mortgageValue, nil
}

func (u *utility) DemortgageValue() (types.Money, error) {
	return u.mortgageValue + u.mortgageValue/10, nil
}
