package gamecore

import (
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

type wallet struct {
	money types.Money
}

type properties struct {
	property_array []tile.Tile
}

type player struct {
	id     types.ID
	wallet wallet
}
