# Environment, Tooling & Test Catalog

This document details the environment configuration, CLI commands, and test suite map for LandGrab.

---

## 1. Environment & CLI Setup

* **Operating System**: Linux x86_64
* **Go Compiler Location**: `/snap/bin/go`  
  *(Note: Standard sandbox PATH does not always include `/snap/bin`. Always invoke `/snap/bin/go` directly when executing test or build commands in terminal tools).*
* **Go Version**: `go1.27.1 linux/amd64`

### Standard Commands

```bash
# Run all unit and integration tests across the workspace
/snap/bin/go test -v ./...

# Run all tests without cache
/snap/bin/go test -count=1 ./...

# Run tests for a specific package
/snap/bin/go test -v ./gamecore/movement/...
/snap/bin/go test -v ./gamecore/decision/...
/snap/bin/go test -v ./gamecore/state/...
/snap/bin/go test -v ./gamecore/pipeline/...
/snap/bin/go test -v ./gamecore/tile/...

# Run a single specific test function
/snap/bin/go test -v ./gamecore/decision -run TestBuildHouseMonopolyAndEvenBuildingRule
```

---

## 2. Test Catalog & Coverage Map

| Test Suite | Package | What It Verifies |
| :--- | :--- | :--- |
| `TestBoardTileCount` | `tile` | Verifies board contains exactly 40 non-nil tiles. |
| `TestBoardCircularWrapping` | `tile` | Circular indexing ($40 \to 0, 41 \to 1, -1 \to 39$). |
| `TestBoardFlyweightIdentity` | `tile` | Pointer identity of repeated `TileAt` calls. |
| `TestBoardTileByID` | `tile` | 0-based lookups for IDs $0..39$ and rejection of out-of-bounds IDs. |
| `TestPropertyTile` | `tile` | Base rent, 1–4 house rent, hotel rent, mortgage values. |
| `TestRailwayTile` | `tile` | 1–4 railway rent scaling ($25 \to $50 \to $100 \to $200). |
| `TestUtilityTile` | `tile` | $4\times$ and $10\times$ dice multipliers. |
| `TestWallet` | `state` | Cash addition, deduction, overdraft protection. |
| `TestPlayerAndProperties` | `state` | Adding, finding, and removing `PropertyHolder` items; counting sets. |
| `TestPropertyHolderRentCalculation`| `state` | Rent dispatch via flyweight; $0 rent when mortgaged. |
| `TestRailwayAndUtilityRentCalculation` | `state` | Railway and utility rent dispatch via `PropertyHolder`. |
| `TestTotalInColorGroup` | `state` | Verifies monopoly set sizes for all 8 colors. |
| `TestGameState` | `state` | Initial state, $O(1)$ `GetPropertyOwner`, multi-player turn cycles. |
| `TestMoveChainExecutionAndSkipping` | `pipeline` | Linear sequencing and `ShouldExecute == false` skipping. |
| `TestDecisionChainExecutionAndSkipping` | `pipeline` | Decision chain sequencing and conditional node execution. |
| `TestChainErrorPropagation` | `pipeline` | Verifies errors halt chain execution and return to caller. |
| `TestMovementToUnownedProperty` | `movement` | Pausing at `PhaseAwaitingDecision` on unowned property. |
| `TestMovementRentPayment` | `movement` | Automatic rent calculation and cash transfer to landlord. |
| `TestSelfOwnedPropertyLanding` | `movement` | Free landing on self-owned tiles. |
| `TestMortgagedPropertyLanding` | `movement` | $0 rent on mortgaged properties. |
| `TestRailwayLandingRent` | `movement` | End-to-end railway landing and rent transfer. |
| `TestUtilityLandingRent` | `movement` | End-to-end utility landing and dice-multiplier rent transfer. |
| `TestBankruptcyOnRent` | `movement` | Liquidation and bankruptcy when rent exceeds wallet. |
| `TestMovementPassingGo` | `movement` | Awards $200 salary when passing or landing on Go. |
| `TestMovementTaxDeductionAndBankruptcy`| `movement` | Income Tax ($200), Wealth Tax ($100), and tax insolvency. |
| `TestGoToJailTile` | `movement` | Arrest when landing on tile 30. |
| `TestThreeDoublesSendToJail` | `movement` | 3 consecutive doubles speeding penalty. |
| `TestJailReleaseOnDoubles` | `movement` | Doubles release from jail without extra roll. |
| `TestJailTurnProgressionAndForcedBail` | `movement` | Turn 1–2 wait, turn 3 forced $50 bail vs bankruptcy. |
| `TestDoublesGrantsExtraTurn` | `movement` | Retaining turn on doubles. |
| `TestBuyDecision` & `Errors` | `decision` | Property purchase, inventory addition, and fund checks. |
| `TestDeclineBuyDecision` | `decision` | Declining property purchase and turn advancement. |
| `TestBuildHouseMonopolyAndEvenBuildingRule` | `decision` | Monopoly enforcement and even-building rule ($\le 1$ house diff). |
| `TestBuildHouseErrors` | `decision` | Unowned, non-buildable, mortgaged in group, max capacity, funds. |
| `TestBuildHotelAndSell` | `decision` | Hotel upgrade (4H $\to$ 1 Hotel) and even-selling with 50% refund. |
| `TestMortgageAndDemortgage` | `decision` | No buildings rule, mortgage credit, 10% demortgage fee. |
| `TestPayBail` | `decision` | $50 voluntary jail exit. |
| `TestTradeDecision` & `Errors` | `decision` | Peer-to-peer property and cash exchange. |
| `TestFullGameTurnLoop` | `gamecore` | End-to-end integration test simulating a multi-turn match. |
