package state

import "github.com/LandGrab/gamecore/types"

type GamePhase int

const (
	PhaseReadyToRoll GamePhase = iota
	PhaseAwaitingDecision
	PhaseTurnEnded
	PhaseGameOver
)

type DecisionType int

const (
	DecisionNone DecisionType = iota
	DecisionBuy
	DecisionDeclineBuy
	DecisionBuildHouse
	DecisionSellHouse
	DecisionMortgage
	DecisionDemortgage
	DecisionTrade
	DecisionPayBail
)

// NoOwnerID represents an unowned tile or unassigned entity.
const NoOwnerID types.ID = ^types.ID(0)
