package tile

import (
	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

type Color uint8

const (
	ColorNone Color = iota
	ColorBrown
	ColorLightBlue
	ColorPink
	ColorOrange
	ColorRed
	ColorYellow
	ColorGreen
	ColorDarkBlue
)

type property struct {
	id            types.ID
	name          string
	color         Color
	cost          types.Money
	rent          types.Money
	rent1House    types.Money
	rent2House    types.Money
	rent3House    types.Money
	rent4House    types.Money
	rentHotel     types.Money
	costHouse     types.Money
	costHotel     types.Money
	mortgageValue types.Money
}

type PropertyConfig struct {
	ID            types.ID
	Name          string
	Color         Color
	Cost          types.Money
	Rent          types.Money
	Rent1House    types.Money
	Rent2House    types.Money
	Rent3House    types.Money
	Rent4House    types.Money
	RentHotel     types.Money
	CostHouse     types.Money
	CostHotel     types.Money
	MortgageValue types.Money
}

func NewProperty(cfg PropertyConfig) Tile {
	return &property{
		id:            cfg.ID,
		name:          cfg.Name,
		color:         cfg.Color,
		cost:          cfg.Cost,
		rent:          cfg.Rent,
		rent1House:    cfg.Rent1House,
		rent2House:    cfg.Rent2House,
		rent3House:    cfg.Rent3House,
		rent4House:    cfg.Rent4House,
		rentHotel:     cfg.RentHotel,
		costHouse:     cfg.CostHouse,
		costHotel:     cfg.CostHotel,
		mortgageValue: cfg.MortgageValue,
	}
}

func (p *property) ID() types.ID {
	return p.id
}

func (p *property) Name() string {
	return p.name
}

func (p *property) Color() Color {
	return p.color
}

func (p *property) Cost() (types.Money, error) {
	return p.cost, nil
}

func (p *property) Rent(info *ExtraInfo) (types.Money, error) {
	if info == nil {
		return 0, errors.ErrInvalidExtraInfo
	}
	if info.Hotel {
		return p.rentHotel, nil
	}
	switch info.Houses {
	case 0:
		return p.rent, nil
	case 1:
		return p.rent1House, nil
	case 2:
		return p.rent2House, nil
	case 3:
		return p.rent3House, nil
	case 4:
		return p.rent4House, nil
	default:
		return 0, errors.ErrInvalidExtraInfo
	}
}

func (p *property) CostHouse() (types.Money, error) {
	return p.costHouse, nil
}

func (p *property) CostHotel() (types.Money, error) {
	return p.costHotel, nil
}

func (p *property) Buildable() bool {
	return true
}

func (p *property) MortgageValue() (types.Money, error) {
	return p.mortgageValue, nil
}

func (p *property) DemortgageValue() (types.Money, error) {
	return p.mortgageValue + p.mortgageValue/10, nil
}
