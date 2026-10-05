# Detailed Package & File Inventory

This file catalogues every directory and source file in `LandGrab`, detailing their symbols, responsibilities, and import relationships.

---

## 1. `gamecore/types`
**Package**: `types` | **Path**: `gamecore/types/types.go`

Defines primitive type aliases to ensure type safety across currency, identification, and dice rolling.
* **`Money uint`**: Currency representation (prevents negative values at the type level).
* **`ID uint`**: Identifier for players and tiles ($0..39$ for tiles).
* **`Dice uint`**: Single die value ($1..6$).

---

## 2. `gamecore/errors`
**Package**: `errors` | **Path**: `gamecore/errors/errors.go`

Defines sentinel errors for tile operations:
* `ErrTileNotFound`: ID exceeds board bounds.
* `ErrNotPurchasable`: Attempting to purchase a non-purchasable tile (e.g. Start, Tax).
* `ErrNotBuildable`: Attempting to build on non-street tiles (e.g. Railway, Utility).
* `ErrNotRentable`: Attempting to collect rent on non-rentable tiles.
* `ErrNoMortgage`: Non-mortgageable tiles.
* `ErrInvalidExtraInfo`: Missing dice or counts when calculating rent.

---

## 3. `gamecore/tile`
**Package**: `tile` | **Path**: `gamecore/tile/`

Defines the static immutable Flyweight board tiles and the 40-tile singleton array.

| File | Primary Symbols | Responsibility |
| :--- | :--- | :--- |
| `tile.go` | `Tile` interface, `ExtraInfo` struct, `Color` enum | Master tile contract and dynamic rent calculation context. |
| `property.go` | `PropertyConfig`, `NewProperty`, `property` | Implements street properties with 6-tier rent tables and house/hotel costs. |
| `railway.go` | `NewRailway`, `railway` | Implements stations ($200 cost, base rent $25, doubling per station owned). |
| `utility.go` | `NewUtility`, `utility` | Implements utilities ($150 cost, rent = $4\times$ or $10\times$ dice sum). |
| `normal.go` | `NewNormal`, `normal` | Implements non-purchasable tiles (Start, Tax, Jail, Fortune, Community Fund). |
| `board.go` | `defaultBoard [40]Tile`, `TileAt(int)`, `TileByID(ID)`, `AllTiles()` | Board singleton with 0-indexed direct lookup. No hashmaps. |
| `board_test.go` | `TestBoardTileCount`, `TestBoardCircularWrapping`, etc. | Verifies circular wrap, flyweight identity, and 0-based lookups. |
| `tile_test.go` | `TestPropertyTile`, `TestRailwayTile`, etc. | Verifies cost, rent calculations, and buildability rules. |

---

## 4. `gamecore/state`
**Package**: `state` | **Path**: `gamecore/state/`

Contains mutable runtime match state. Acts as the single source of truth passed down the execution pipelines.

| File | Primary Symbols | Responsibility |
| :--- | :--- | :--- |
| `enums.go` | `GamePhase`, `DecisionType`, `NoOwnerID` | Turn phase machine, pending decision machine, and unowned sentinel. |
| `property_holder.go`| `PropertyHolder`, `NewPropertyHolder`, `CalculateRent` | Flyweight bridge: bundles static `Tile` with `IsMortgaged` and `ExtraInfo`. Computes rent. |
| `player.go` | `Player`, `Wallet`, `NewPlayer`, helper methods | Player entity, wallet debit/credit, property inventory, color/railway/utility counting. |
| `game.go` | `Game`, `NewGame`, `AdvanceTurn`, `GetPropertyOwner` | Root match struct, $O(1)$ `TileOwnerIndex [40]types.ID`, turn cycling. |
| `state_test.go` | Unit test suite | Tests wallet, player inventory, property holder rent, color counting, and game state. |

---

## 5. `gamecore/pipeline`
**Package**: `pipeline` | **Path**: `gamecore/pipeline/`

Defines the Chain of Responsibility abstraction.

| File | Primary Symbols | Responsibility |
| :--- | :--- | :--- |
| `move_node.go` | `MoveNode` interface, `BaseMoveNode` struct | Defines `ShouldExecute(g) bool`, `Execute(g) error`, `SetNext`, and `ExecuteNext`. |
| `decision_node.go`| `DecisionNode` interface, `BaseDecisionNode` struct| Defines decision chain link with `ShouldExecute(g) bool` and `Execute(g) error`. |
| `builder.go` | `BuildMoveChain`, `BuildDecisionChain` | Variadic helper linking nodes into a sequential chain and returning the head. |
| `pipeline_test.go` | Unit test suite | Tests chain sequencing, `ShouldExecute` skip logic, and error propagation. |

---

## 6. `gamecore/movement`
**Package**: `movement` | **Path**: `gamecore/movement/`

Implements the 9-node sequential movement pipeline driving player turns.

| File | Node Implemented | Responsibility |
| :--- | :--- | :--- |
| `dice_node.go` | `DiceRollNode` | Dice rolling, consecutive doubles counting, 3-doubles speeding arrest to jail. |
| `jail_node.go` | `JailCheckNode` | Doubles escape without extra turn, turns 1–2 progression, turn 3 mandatory $50 bail. |
| `position_node.go` | `PositionUpdateNode` | 40-tile circular wrapping and $200 Go salary credit. |
| `special_nodes.go` | `GoToJailNode`, `TaxNode` | Tile #30 arrest to pos 10; Income Tax (pos 4 $200) & Wealth Tax (pos 38 $100). |
| `landing_nodes.go` | `PropertyLandingNode`, `RailwayLandingNode`, `UtilityLandingNode` | Ownership resolution, purchase prompt pause (`PhaseAwaitingDecision`), rent calculation/transfer. |
| `turn_end_node.go` | `TurnEndNode` | Grants extra roll on doubles (`PhaseReadyToRoll`) or advances turn (`AdvanceTurn`). |
| `manager.go` | `MoveManager`, `NewMoveManager(roller)` | Assembles the 9 nodes and executes the movement pipeline on `*state.Game`. |
| `movement_test.go` | Integration test suite | Tests all movement scenarios, taxes, jail escape, rent, and bankruptcies. |

---

## 7. `gamecore/decision`
**Package**: `decision` | **Path**: `gamecore/decision/`

Implements the 9-node sequential decision pipeline processing player choices and business rules.

| File | Node / Symbols | Responsibility |
| :--- | :--- | :--- |
| `turn_validation_node.go`| `TurnValidationNode`, Sentinel Errors | Sanity-checks active player; declares all domain errors (`ErrMonopolyRequired`, etc.). |
| `buy_node.go` | `BuyDecisionNode`, `DeclineBuyDecisionNode` | Purchases property, adds `PropertyHolder`, updates `TileOwnerIndex`, resumes turn. |
| `build_node.go` | `BuildHouseDecisionNode` | Enforces Monopoly Rule, No-Mortgage in Group Rule, Even-Building Rule, Hotel upgrade. |
| `sell_node.go` | `SellHouseDecisionNode` | Enforces Even-Selling Rule, downgrades hotel to 4 houses, refunds 50% build cost. |
| `mortgage_node.go` | `MortgageDecisionNode`, `DemortgageDecisionNode` | Enforces No-Buildings in Group Rule, grants mortgage value, charges 10% interest to demortgage. |
| `trade_node.go` | `TradeDecisionNode`, `TradeOffer` | Peer-to-peer property transfer and cash settlement between players. |
| `bail_node.go` | `PayBailDecisionNode` | Voluntary $50 bail payment to exit jail immediately. |
| `manager.go` | `DecisionManager`, `NewDecisionManager` | Assembles the 9 decision nodes and executes decision logic on `*state.Game`. |
| `decision_test.go` | Integration test suite | Comprehensive validation of all decision rules, limits, and error paths. |

---

## 8. Root & Integration Files

| File | Responsibility |
| :--- | :--- |
| `main.go` | Application entry point. |
| `go.mod` | Go module descriptor (`module github.com/LandGrab`, Go 1.27.1). |
| `.gitignore` | Ignores `llm_memory/`. |
| `gamecore/integration_test.go` | End-to-end multi-turn game simulation testing buying, monopolies, house building, and rent transfers. |
