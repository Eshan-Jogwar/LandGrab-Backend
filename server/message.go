package server

import "github.com/LandGrab/gamecore/types"

// Client action types sent from client -> server over WebSocket.
const (
	ActionRollDice     = "roll_dice"
	ActionBuyProperty  = "buy_property"
	ActionDeclineBuy   = "decline_buy"
	ActionBuildHouse   = "build_house"
	ActionSellHouse    = "sell_house"
	ActionMortgage     = "mortgage"
	ActionDemortgage   = "demortgage"
	ActionTrade        = "trade"
	ActionPayBail      = "pay_bail"
)

// Server message types sent from server -> client over WebSocket.
const (
	MsgConnected          = "connected"
	MsgQueued             = "queued"
	MsgGameStart          = "game_start"
	MsgGameState          = "game_state"
	MsgYourTurn           = "your_turn"
	MsgActionResult       = "action_result"
	MsgError              = "error"
	MsgPlayerDisconnected = "player_disconnected"
	MsgGameOver           = "game_over"
)

// ClientMessage represents an incoming JSON action payload from a player's WebSocket.
type ClientMessage struct {
	Type       string      `json:"type"`
	TileID     *types.ID   `json:"tile_id,omitempty"`
	ToPlayerID *types.ID   `json:"to_player_id,omitempty"`
	Cash       *types.Money `json:"cash,omitempty"`
}

// ServerMessage represents an outgoing JSON envelope broadcast to one or all players.
type ServerMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ConnectedPayload is sent upon successful WebSocket handshake.
type ConnectedPayload struct {
	TempID string `json:"temp_id"`
}

// QueuedPayload informs the client of their position and current queue state.
type QueuedPayload struct {
	TempID       string `json:"temp_id"`
	Position     int    `json:"position"`
	TotalQueued  int    `json:"total_queued"`
	RequiredSize int    `json:"required_size"`
}

// GameStartPayload informs players that a match has been created.
type GameStartPayload struct {
	GameID        string             `json:"game_id"`
	AssignedID    types.ID           `json:"assigned_id"`
	TempID        string             `json:"temp_id"`
	TotalPlayers  int                `json:"total_players"`
	StartingMoney types.Money        `json:"starting_money"`
	Players       []PlayerInfoSummary `json:"players"`
}

// PlayerInfoSummary is lightweight metadata sent during game start.
type PlayerInfoSummary struct {
	ID     types.ID `json:"id"`
	TempID string   `json:"temp_id"`
}

// PropertySnapshot provides serializable status for an individual property tile.
type PropertySnapshot struct {
	TileID      types.ID `json:"tile_id"`
	TileName    string   `json:"tile_name"`
	IsMortgaged bool     `json:"is_mortgaged"`
	Houses      int      `json:"houses"`
	HasHotel    bool     `json:"has_hotel"`
}

// PlayerSnapshot contains full state details for a player in a game.
type PlayerSnapshot struct {
	ID         types.ID           `json:"id"`
	TempID     string             `json:"temp_id"`
	Balance    types.Money        `json:"balance"`
	Position   int                `json:"position"`
	InJail     bool               `json:"in_jail"`
	IsBankrupt bool               `json:"is_bankrupt"`
	Properties []PropertySnapshot `json:"properties"`
}

// GameStatePayload provides the full snapshot of match state after each turn/action.
type GameStatePayload struct {
	GameID             string           `json:"game_id"`
	CurrentPlayerID    types.ID         `json:"current_player_id"`
	Phase              string           `json:"phase"`
	DiceRoll           [2]types.Dice    `json:"dice_roll"`
	ConsecutiveDoubles int              `json:"consecutive_doubles"`
	CanRollAgain       bool             `json:"can_roll_again"`
	PendingDecision    string           `json:"pending_decision,omitempty"`
	DecisionTargetTile types.ID         `json:"decision_target_tile,omitempty"`
	Players            []PlayerSnapshot `json:"players"`
}

// YourTurnPayload prompts the active player with allowable actions.
type YourTurnPayload struct {
	CurrentPlayerID types.ID `json:"current_player_id"`
	Phase           string   `json:"phase"`
	PendingDecision string   `json:"pending_decision,omitempty"`
	TimeoutSeconds  int      `json:"timeout_seconds"`
}

// ActionResultPayload informs the player about the outcome of their attempted move.
type ActionResultPayload struct {
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// GameOverPayload announces the winner and final standings.
type GameOverPayload struct {
	GameID       string   `json:"game_id"`
	WinnerID     types.ID `json:"winner_id"`
	WinnerTempID string   `json:"winner_temp_id"`
}
