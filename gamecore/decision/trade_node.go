package decision

import (
	"errors"

	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/types"
)

var ErrTradeInvalid = errors.New("invalid trade configuration")

type TradeOffer struct {
	FromPlayerID types.ID
	ToPlayerID   types.ID
	TileID       types.ID
	Cash         types.Money
}

type TradeDecisionNode struct {
	pipeline.BaseDecisionNode
	PendingTrade *TradeOffer
}

func (n *TradeDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionTrade
}

func (n *TradeDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	if n.PendingTrade == nil {
		g.PendingDecision = state.DecisionNone
		return ErrTradeInvalid
	}

	trade := n.PendingTrade
	var fromPlayer, toPlayer *state.Player
	for _, p := range g.Players {
		if p.ID == trade.FromPlayerID {
			fromPlayer = p
		}
		if p.ID == trade.ToPlayerID {
			toPlayer = p
		}
	}

	if fromPlayer == nil || toPlayer == nil {
		g.PendingDecision = state.DecisionNone
		return ErrTradeInvalid
	}

	ph, ok := fromPlayer.RemoveProperty(trade.TileID)
	if !ok {
		g.PendingDecision = state.DecisionNone
		return ErrPropertyNotOwned
	}

	if trade.Cash > 0 {
		if !toPlayer.Wallet.Deduct(trade.Cash) {
			fromPlayer.Properties = append(fromPlayer.Properties, ph)
			g.PendingDecision = state.DecisionNone
			return ErrInsufficientFunds
		}
		fromPlayer.Wallet.Add(trade.Cash)
	}

	toPlayer.Properties = append(toPlayer.Properties, ph)
	g.TileOwnerIndex[trade.TileID] = toPlayer.ID
	g.PendingDecision = state.DecisionNone

	return n.ExecuteNext(g)
}
