package decision

import (
	"errors"

	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

var (
	ErrNoCurrentPlayer     = errors.New("game has no active player")
	ErrInvalidDecision     = errors.New("invalid decision for current game state")
	ErrInsufficientFunds   = errors.New("player has insufficient funds")
	ErrPropertyNotOwned    = errors.New("player does not own the target property")
	ErrMonopolyRequired    = errors.New("cannot build without owning all properties of the color group")
	ErrMortgagedInGroup    = errors.New("cannot build when a property in the color group is mortgaged")
	ErrEvenBuildingRule    = errors.New("must build evenly across properties of the color group")
	ErrMaxHousesReached    = errors.New("cannot build further: hotel already constructed")
	ErrNoHousesToSell      = errors.New("no houses to sell on this property")
	ErrEvenSellingRule     = errors.New("must sell evenly across properties of the color group")
	ErrAlreadyMortgaged    = errors.New("property is already mortgaged")
	ErrNotMortgaged        = errors.New("property is not mortgaged")
	ErrBuildingsInGroup    = errors.New("must sell all houses and hotels in color group before mortgaging")
	ErrPlayerNotInJail     = errors.New("player is not in jail")
)

type TurnValidationNode struct {
	pipeline.BaseDecisionNode
}

func (n *TurnValidationNode) ShouldExecute(g *state.Game) bool {
	return true
}

func (n *TurnValidationNode) Execute(g *state.Game) error {
	if g.CurrentPlayer() == nil {
		return ErrNoCurrentPlayer
	}
	return n.ExecuteNext(g)
}
