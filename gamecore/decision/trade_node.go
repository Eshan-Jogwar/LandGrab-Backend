package decision

import (
	"errors"

	"github.com/LandGrab/gamecore/pipeline"
	"github.com/LandGrab/gamecore/state"
)

var ErrTradeInvalid = errors.New("invalid trade configuration")

// TradeOffer aliases state.TradeOffer for package-level convenience and backward compatibility.
type TradeOffer = state.TradeOffer

type TradeDecisionNode struct {
	pipeline.BaseDecisionNode
}

func (n *TradeDecisionNode) ShouldExecute(g *state.Game) bool {
	return g.PendingDecision == state.DecisionTrade
}

func (n *TradeDecisionNode) Execute(g *state.Game) error {
	if !n.ShouldExecute(g) {
		return n.ExecuteNext(g)
	}

	if g.PendingTrade == nil {
		g.PendingDecision = state.DecisionNone
		return ErrTradeInvalid
	}

	trade := g.PendingTrade
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
		g.PendingTrade = nil
		g.PendingDecision = state.DecisionNone
		return ErrTradeInvalid
	}

	ph, ok := fromPlayer.RemoveProperty(trade.TileID)
	if !ok {
		g.PendingTrade = nil
		g.PendingDecision = state.DecisionNone
		return ErrPropertyNotOwned
	}

	if trade.Cash > 0 {
		if !toPlayer.Wallet.Deduct(trade.Cash) {
			fromPlayer.Properties = append(fromPlayer.Properties, ph)
			g.PendingTrade = nil
			g.PendingDecision = state.DecisionNone
			return ErrInsufficientFunds
		}
		fromPlayer.Wallet.Add(trade.Cash)
	}

	toPlayer.Properties = append(toPlayer.Properties, ph)
	g.TileOwnerIndex[trade.TileID] = toPlayer.ID
	g.PendingDecision = state.DecisionNone
	g.PendingTrade = nil

	return n.ExecuteNext(g)
}
