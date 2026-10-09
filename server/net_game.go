package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/LandGrab/gamecore/decision"
	"github.com/LandGrab/gamecore/movement"
	"github.com/LandGrab/gamecore/state"
	"github.com/LandGrab/gamecore/tile"
	"github.com/LandGrab/gamecore/types"
)

const (
	defaultTurnTimeout = 45 * time.Second
)

// PlayerAction represents an action sent from a specific connected player.
type PlayerAction struct {
	SenderTempID string
	Message      ClientMessage
}

// NetGame wraps the gamecore Game and manages the real-time session event loop.
type NetGame struct {
	ID          string
	Game        *state.Game
	Players     [4]*NetPlayer
	PlayerIndex map[types.ID]*NetPlayer
	TempIDIndex map[string]*NetPlayer

	actions   chan PlayerAction
	done      chan struct{}
	turnTimer *time.Timer
	turnLimit time.Duration
	mu        sync.Mutex
	CreatedAt time.Time

	onFinish func(gameID string)
}

// NewNetGame constructs a NetGame session with 4 players.
func NewNetGame(id string, players [4]*NetPlayer) *NetGame {
	corePlayers := make([]*state.Player, 4)
	playerIndex := make(map[types.ID]*NetPlayer)
	tempIDIndex := make(map[string]*NetPlayer)

	for i := 0; i < 4; i++ {
		pID := types.ID(i + 1)
		cp := state.NewPlayer(pID, 1500)
		corePlayers[i] = cp

		netPlayer := players[i]
		netPlayer.Player = cp
		netPlayer.GameID = id

		playerIndex[pID] = netPlayer
		tempIDIndex[netPlayer.TempID] = netPlayer
	}

	game := state.NewGame(corePlayers)

	return &NetGame{
		ID:          id,
		Game:        game,
		Players:     players,
		PlayerIndex: playerIndex,
		TempIDIndex: tempIDIndex,
		actions:     make(chan PlayerAction, 64),
		done:        make(chan struct{}),
		turnLimit:   defaultTurnTimeout,
		CreatedAt:   time.Now(),
	}
}

// Run executes the session actor loop in a dedicated goroutine.
func (ng *NetGame) Run() {
	// Bind action handlers to all players
	for _, p := range ng.Players {
		p.SetActionHandler(func(action PlayerAction) {
			select {
			case ng.actions <- action:
			case <-ng.done:
			}
		})
		p.SetDisconnectHandler(func(tempID string) {
			select {
			case ng.actions <- PlayerAction{
				SenderTempID: tempID,
				Message: ClientMessage{
					Type: "internal_disconnect",
				},
			}:
			case <-ng.done:
			}
		})
	}

	// Send initial game_start notification to all players
	ng.broadcastGameStart()
	ng.broadcastGameState()
	ng.resetTurnTimer()
	ng.notifyCurrentPlayer()

	defer ng.cleanup()

	for {
		select {
		case action := <-ng.actions:
			ng.handleAction(action)
		case <-ng.turnTimer.C:
			ng.handleTurnTimeout()
		case <-ng.done:
			return
		}
	}
}

// Stop gracefully stops the session loop.
func (ng *NetGame) Stop() {
	select {
	case <-ng.done:
		return
	default:
		close(ng.done)
	}
}

func (ng *NetGame) cleanup() {
	if ng.turnTimer != nil {
		ng.turnTimer.Stop()
	}
	if ng.onFinish != nil {
		ng.onFinish(ng.ID)
	}
}

func (ng *NetGame) broadcastGameStart() {
	summaries := make([]PlayerInfoSummary, 4)
	for i, np := range ng.Players {
		summaries[i] = PlayerInfoSummary{
			ID:     np.Player.ID,
			TempID: np.TempID,
		}
	}

	for _, np := range ng.Players {
		np.SendJSON(ServerMessage{
			Type: MsgGameStart,
			Payload: GameStartPayload{
				GameID:        ng.ID,
				AssignedID:    np.Player.ID,
				TempID:        np.TempID,
				TotalPlayers:  4,
				StartingMoney: 1500,
				Players:       summaries,
			},
		})
	}
}

func (ng *NetGame) broadcastGameState() {
	snapshot := ng.CreateSnapshot()
	msg := ServerMessage{
		Type:    MsgGameState,
		Payload: snapshot,
	}
	for _, np := range ng.Players {
		np.SendJSON(msg)
	}
}

func (ng *NetGame) notifyCurrentPlayer() {
	cp := ng.Game.CurrentPlayer()
	if cp == nil {
		return
	}

	currNetPlayer, ok := ng.PlayerIndex[cp.ID]
	if !ok {
		return
	}

	phaseStr := ng.phaseString(ng.Game.Phase)
	decisionStr := ng.decisionString(ng.Game.PendingDecision)

	currNetPlayer.SendJSON(ServerMessage{
		Type: MsgYourTurn,
		Payload: YourTurnPayload{
			CurrentPlayerID: cp.ID,
			Phase:           phaseStr,
			PendingDecision: decisionStr,
			TimeoutSeconds:  int(ng.turnLimit.Seconds()),
		},
	})
}

func (ng *NetGame) resetTurnTimer() {
	if ng.turnTimer != nil {
		ng.turnTimer.Stop()
	}
	ng.turnTimer = time.NewTimer(ng.turnLimit)
}

func (ng *NetGame) handleAction(action PlayerAction) {
	if action.Message.Type == "internal_disconnect" {
		ng.handleDisconnect(action.SenderTempID)
		return
	}

	player, ok := ng.TempIDIndex[action.SenderTempID]
	if !ok {
		return
	}

	cp := ng.Game.CurrentPlayer()
	if cp == nil {
		player.SendJSON(ServerMessage{Type: MsgError, Error: "Game has no active player"})
		return
	}

	// 1. Validation: Verify sender is the current player
	if player.Player.ID != cp.ID {
		player.SendJSON(ServerMessage{
			Type:  MsgError,
			Error: fmt.Sprintf("It is currently Player %d's turn", cp.ID),
		})
		return
	}

	// 2. Dispatch action
	var err error
	switch action.Message.Type {
	case ActionRollDice:
		if ng.Game.Phase != state.PhaseReadyToRoll {
			err = fmt.Errorf("cannot roll dice in phase %s", ng.phaseString(ng.Game.Phase))
			break
		}
		err = movement.Execute(ng.Game)

	case ActionBuyProperty:
		if ng.Game.Phase != state.PhaseAwaitingDecision || ng.Game.PendingDecision != state.DecisionBuy {
			err = fmt.Errorf("no property purchase decision pending")
			break
		}
		err = decision.Execute(ng.Game)

	case ActionDeclineBuy:
		if ng.Game.Phase != state.PhaseAwaitingDecision || ng.Game.PendingDecision != state.DecisionBuy {
			err = fmt.Errorf("no property purchase decision pending")
			break
		}
		ng.Game.PendingDecision = state.DecisionDeclineBuy
		err = decision.Execute(ng.Game)

	case ActionBuildHouse:
		if action.Message.TileID == nil {
			err = fmt.Errorf("tile_id is required to build a house")
			break
		}
		ng.Game.PendingDecision = state.DecisionBuildHouse
		ng.Game.DecisionTargetTile = *action.Message.TileID
		err = decision.Execute(ng.Game)

	case ActionSellHouse:
		if action.Message.TileID == nil {
			err = fmt.Errorf("tile_id is required to sell a house")
			break
		}
		ng.Game.PendingDecision = state.DecisionSellHouse
		ng.Game.DecisionTargetTile = *action.Message.TileID
		err = decision.Execute(ng.Game)

	case ActionMortgage:
		if action.Message.TileID == nil {
			err = fmt.Errorf("tile_id is required to mortgage property")
			break
		}
		ng.Game.PendingDecision = state.DecisionMortgage
		ng.Game.DecisionTargetTile = *action.Message.TileID
		err = decision.Execute(ng.Game)

	case ActionDemortgage:
		if action.Message.TileID == nil {
			err = fmt.Errorf("tile_id is required to demortgage property")
			break
		}
		ng.Game.PendingDecision = state.DecisionDemortgage
		ng.Game.DecisionTargetTile = *action.Message.TileID
		err = decision.Execute(ng.Game)

	case ActionTrade:
		if action.Message.ToPlayerID == nil || action.Message.TileID == nil {
			err = fmt.Errorf("to_player_id and tile_id are required for trading")
			break
		}
		cash := types.Money(0)
		if action.Message.Cash != nil {
			cash = *action.Message.Cash
		}
		offer := &state.TradeOffer{
			FromPlayerID: cp.ID,
			ToPlayerID:   *action.Message.ToPlayerID,
			TileID:       *action.Message.TileID,
			Cash:         cash,
		}
		err = decision.ExecuteTrade(ng.Game, offer)

	case ActionPayBail:
		if !cp.InJail {
			err = fmt.Errorf("player is not in jail")
			break
		}
		ng.Game.PendingDecision = state.DecisionPayBail
		err = decision.Execute(ng.Game)

	default:
		err = fmt.Errorf("unknown action type: %s", action.Message.Type)
	}

	if err != nil {
		player.SendJSON(ServerMessage{
			Type:  MsgError,
			Error: err.Error(),
		})
		return
	}

	player.SendJSON(ServerMessage{
		Type: MsgActionResult,
		Payload: ActionResultPayload{
			Action:  action.Message.Type,
			Success: true,
		},
	})

	// Check win condition
	if ng.checkGameOver() {
		return
	}

	// Broadcast new state and reset turn timer
	ng.broadcastGameState()
	ng.resetTurnTimer()
	ng.notifyCurrentPlayer()
}

func (ng *NetGame) handleTurnTimeout() {
	cp := ng.Game.CurrentPlayer()
	if cp == nil {
		return
	}

	// Auto-resolve turn if player goes AFK
	if ng.Game.Phase == state.PhaseReadyToRoll {
		_ = movement.Execute(ng.Game)
	} else if ng.Game.Phase == state.PhaseAwaitingDecision {
		ng.Game.PendingDecision = state.DecisionDeclineBuy
		_ = decision.Execute(ng.Game)
	}

	if ng.checkGameOver() {
		return
	}

	ng.broadcastGameState()
	ng.resetTurnTimer()
	ng.notifyCurrentPlayer()
}

func (ng *NetGame) handleDisconnect(tempID string) {
	np, ok := ng.TempIDIndex[tempID]
	if !ok {
		return
	}

	// Broadcast disconnection notice
	msg := ServerMessage{
		Type: MsgPlayerDisconnected,
		Payload: map[string]interface{}{
			"temp_id":   tempID,
			"player_id": np.Player.ID,
		},
	}
	for _, p := range ng.Players {
		if p.TempID != tempID {
			p.SendJSON(msg)
		}
	}

	// If disconnected player is active, auto-advance
	cp := ng.Game.CurrentPlayer()
	if cp != nil && cp.ID == np.Player.ID {
		ng.handleTurnTimeout()
	}
}

func (ng *NetGame) checkGameOver() bool {
	activeCount := 0
	var lastStanding *NetPlayer

	for _, np := range ng.Players {
		if !np.Player.IsBankrupt {
			activeCount++
			lastStanding = np
		}
	}

	if activeCount <= 1 && lastStanding != nil {
		ng.broadcastGameOver(lastStanding)
		ng.Stop()
		return true
	}
	return false
}

func (ng *NetGame) broadcastGameOver(winner *NetPlayer) {
	msg := ServerMessage{
		Type: MsgGameOver,
		Payload: GameOverPayload{
			GameID:       ng.ID,
			WinnerID:     winner.Player.ID,
			WinnerTempID: winner.TempID,
		},
	}
	for _, np := range ng.Players {
		np.SendJSON(msg)
	}
}

// CreateSnapshot serializes the current Game state into a JSON-friendly struct.
func (ng *NetGame) CreateSnapshot() GameStatePayload {
	cpID := types.ID(0)
	if cp := ng.Game.CurrentPlayer(); cp != nil {
		cpID = cp.ID
	}

	players := make([]PlayerSnapshot, 4)
	for i, np := range ng.Players {
		props := make([]PropertySnapshot, 0, len(np.Player.Properties))
		for _, ph := range np.Player.Properties {
			props = append(props, PropertySnapshot{
				TileID:      ph.Tile.ID(),
				TileName:    ph.Tile.Name(),
				IsMortgaged: ph.IsMortgaged,
				Houses:      ph.ExtraInfo.Houses,
				HasHotel:    ph.ExtraInfo.Hotel,
			})
		}

		players[i] = PlayerSnapshot{
			ID:         np.Player.ID,
			TempID:     np.TempID,
			Balance:    np.Player.Wallet.Balance(),
			Position:   np.Player.Position,
			InJail:     np.Player.InJail,
			IsBankrupt: np.Player.IsBankrupt,
			Properties: props,
		}
	}

	return GameStatePayload{
		GameID:             ng.ID,
		CurrentPlayerID:    cpID,
		Phase:              ng.phaseString(ng.Game.Phase),
		DiceRoll:           ng.Game.DiceRoll,
		ConsecutiveDoubles: ng.Game.ConsecutiveDoubles,
		CanRollAgain:       ng.Game.CanRollAgain,
		PendingDecision:    ng.decisionString(ng.Game.PendingDecision),
		DecisionTargetTile: ng.Game.DecisionTargetTile,
		Players:            players,
	}
}

func (ng *NetGame) phaseString(phase state.GamePhase) string {
	switch phase {
	case state.PhaseReadyToRoll:
		return "ready_to_roll"
	case state.PhaseAwaitingDecision:
		return "awaiting_decision"
	case state.PhaseTurnEnded:
		return "turn_ended"
	case state.PhaseGameOver:
		return "game_over"
	default:
		return "unknown"
	}
}

func (ng *NetGame) decisionString(decision state.DecisionType) string {
	switch decision {
	case state.DecisionBuy:
		return "buy"
	case state.DecisionDeclineBuy:
		return "decline_buy"
	case state.DecisionBuildHouse:
		return "build_house"
	case state.DecisionSellHouse:
		return "sell_house"
	case state.DecisionMortgage:
		return "mortgage"
	case state.DecisionDemortgage:
		return "demortgage"
	case state.DecisionTrade:
		return "trade"
	case state.DecisionPayBail:
		return "pay_bail"
	default:
		return "none"
	}
}

// CurrentTileName returns human-readable name of current tile.
func (ng *NetGame) CurrentTileName() string {
	t := ng.Game.CurrentTile()
	if t == nil {
		return tile.TileAt(0).Name()
	}
	return t.Name()
}
