package tile

import (
	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

type railway struct {
	id            types.ID
	name          string
	cost          types.Money
	baseRent      types.Money
	mortgageValue types.Money
}

func NewRailway(id types.ID, name string, cost types.Money, baseRent types.Money, mortgageValue types.Money) Tile {
	return &railway{
		id:            id,
		name:          name,
		cost:          cost,
		baseRent:      baseRent,
		mortgageValue: mortgageValue,
	}
}

func (r *railway) ID() types.ID {
	return r.id
}

func (r *railway) Color() Color {
	return ColorNone
}

func (r *railway) Name() string {
	return r.name
}

func (r *railway) Cost() (types.Money, error) {
	return r.cost, nil
}

func (r *railway) Rent(info *ExtraInfo) (types.Money, error) {
	if info == nil || info.RailwaysOwned <= 0 {
		return 0, errors.ErrInvalidExtraInfo
	}
	multiplier := types.Money(1 << (info.RailwaysOwned - 1))
	return r.baseRent * multiplier, nil
}

func (r *railway) CostHouse() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (r *railway) CostHotel() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (r *railway) Buildable() bool {
	return false
}

func (r *railway) MortgageValue() (types.Money, error) {
	return r.mortgageValue, nil
}

func (r *railway) DemortgageValue() (types.Money, error) {
	return r.mortgageValue + r.mortgageValue/10, nil
}
