package tile

import (
	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

type normal struct {
	id   types.ID
	name string
}

func NewNormal(id types.ID, name string) Tile {
	return &normal{
		id:   id,
		name: name,
	}
}

func (n *normal) ID() types.ID {
	return n.id
}

func (n *normal) Color() Color {
	return ColorNone
}

func (n *normal) Name() string {
	return n.name
}

func (n *normal) Cost() (types.Money, error) {
	return 0, errors.ErrNotPurchasable
}

func (n *normal) Rent(*ExtraInfo) (types.Money, error) {
	return 0, errors.ErrNotRentable
}

func (n *normal) CostHouse() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (n *normal) CostHotel() (types.Money, error) {
	return 0, errors.ErrNotBuildable
}

func (n *normal) Buildable() bool {
	return false
}

func (n *normal) MortgageValue() (types.Money, error) {
	return 0, errors.ErrNoMortgage
}

func (n *normal) DemortgageValue() (types.Money, error) {
	return 0, errors.ErrNoMortgage
}
