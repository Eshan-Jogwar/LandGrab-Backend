# LLM Memory & Master Context Index

> **LLM Instructions for Future Sessions**:  
> Read this directory first when starting any session on `LandGrab`. It contains the complete architectural map, file inventory, domain rules, and testing commands, eliminating the need for exploratory file searches (`find`, `grep`, `ls`).

---

## 1. Quick Repository Map

```
LandGrab/
├── .gitignore                     # Ignores llm_memory/
├── go.mod                         # Module: github.com/LandGrab (Go 1.27.1)
├── main.go                        # Root entry point
├── docs/                          # Human-readable architectural documentation
│   ├── README.md, structs.md, movenodes.md, decisionnodes.md
├── llm_memory/                    # Fast-lookup LLM context cache
│   ├── README.md                  # This file (master navigation index)
│   ├── 01_packages_and_files.md   # Exhaustive file-by-file symbol & responsibility index
│   ├── 02_architecture_patterns.md# Flyweight, Chain of Responsibility & state decisions
│   ├── 03_domain_rules.md         # 40-tile board map, monopoly sets, rent formulas, jail/tax rules
│   └── 04_testing_and_tooling.md  # CLI commands (/snap/bin/go) and test suite catalog
└── gamecore/                      # Core game engine
    ├── types/                     # Primitive types (Money, ID, Dice)
    ├── errors/                    # Domain sentinel errors
    ├── tile/                      # 40 immutable flyweight tiles & board array
    ├── state/                     # Mutable game state (Game, Player, Wallet, PropertyHolder)
    ├── pipeline/                  # Chain of Responsibility interfaces & builders
    ├── movement/                  # 9-node linear movement execution pipeline & MoveManager
    ├── decision/                  # 9-node linear decision execution pipeline & DecisionManager
    └── integration_test.go        # Full multi-turn game loop integration test
```

---

## 2. Quick Symbol & File Lookup

| If you need to work on... | Look in this file | Key Symbols |
| :--- | :--- | :--- |
| **Game State & Turn Management** | [`gamecore/state/game.go`](file:///home/eshan/coding/LandGrab/gamecore/state/game.go) | `Game`, `NewGame`, `AdvanceTurn`, `GetPropertyOwner` |
| **Player, Wallet & Ownership** | [`gamecore/state/player.go`](file:///home/eshan/coding/LandGrab/gamecore/state/player.go) | `Player`, `Wallet`, `NewPlayer`, `CountColorGroup` |
| **Property Runtime State & Rent** | [`gamecore/state/property_holder.go`](file:///home/eshan/coding/LandGrab/gamecore/state/property_holder.go) | `PropertyHolder`, `CalculateRent` |
| **Game & Decision Enums** | [`gamecore/state/enums.go`](file:///home/eshan/coding/LandGrab/gamecore/state/enums.go) | `GamePhase`, `DecisionType`, `NoOwnerID` |
| **Movement Pipeline & Execution** | [`gamecore/movement/manager.go`](file:///home/eshan/coding/LandGrab/gamecore/movement/manager.go) | `MoveManager`, `NewMoveManager`, `Execute` |
| **Movement Nodes (Dice, Jail, Movement)** | `gamecore/movement/` | `DiceRollNode`, `JailCheckNode`, `PositionUpdateNode` |
| **Movement Nodes (Taxes & Landings)** | `gamecore/movement/` | `GoToJailNode`, `TaxNode`, `PropertyLandingNode`, etc. |
| **Decision Pipeline & Execution** | [`gamecore/decision/manager.go`](file:///home/eshan/coding/LandGrab/gamecore/decision/manager.go) | `DecisionManager`, `NewDecisionManager`, `Execute` |
| **Decision Nodes (Buy, Build, Mortgage)**| `gamecore/decision/` | `BuyDecisionNode`, `BuildHouseDecisionNode`, `MortgageDecisionNode` |
| **Decision Nodes (Trade, Bail, Sell)** | `gamecore/decision/` | `TradeDecisionNode`, `PayBailDecisionNode`, `SellHouseDecisionNode` |
| **Chain Interfaces & Builders** | `gamecore/pipeline/` | `MoveNode`, `DecisionNode`, `BuildMoveChain` |
| **Static Board & 40 Tiles** | [`gamecore/tile/board.go`](file:///home/eshan/coding/LandGrab/gamecore/tile/board.go) | `TileAt`, `TileByID`, `AllTiles`, `defaultBoard` |
| **Tile Definitions** | `gamecore/tile/` | `Tile`, `PropertyConfig`, `NewProperty`, `NewRailway` |
| **Domain Primitive Types** | [`gamecore/types/types.go`](file:///home/eshan/coding/LandGrab/gamecore/types/types.go) | `Money`, `ID`, `Dice` (all `uint`) |
| **Domain Sentinel Errors** | `gamecore/errors/`, `decision/` | `ErrTileNotFound`, `ErrMonopolyRequired`, etc. |

---

## 3. Essential Environment Invariants

* **Go Executable**: `/snap/bin/go` (not in default sandboxed PATH).
* **Test Command**: `/snap/bin/go test -v ./...`
* **Tile Indexing**: Contiguous **0-based indexing ($0..39$)**. ID 0 is "Start", ID 1 is "Old Harbor Road", ID 39 is "Celestial Heights".
* **Architecture**: Unidirectional DAG: `types` $\to$ `tile` $\to$ `state` $\to$ `pipeline` $\to$ `movement` & `decision`. No circular imports.
