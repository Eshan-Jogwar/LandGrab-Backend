# Domain Rules, Formulas & Constants Cheatsheet

This cheatsheet provides the exact mathematical formulas, game constants, and board configurations in LandGrab.

---

## 1. Complete 40-Tile Board Map

| ID | Name | Category | Color Group | Cost | Base Rent | House/Hotel Cost | Mortgage Value |
| :---: | :--- | :--- | :--- | :---: | :---: | :---: | :---: |
| **0** | Start (Go) | Normal | None | — | — | — | — |
| **1** | Old Harbor Road | Property | Brown | $60 | $2 | $50 | $30 |
| **2** | Community Fund | Normal | None | — | — | — | — |
| **3** | Fisherman's Wharf | Property | Brown | $60 | $4 | $50 | $30 |
| **4** | Income Tax | Normal | None | — | — (Tax: $200) | — | — |
| **5** | Metro Central Station | Railway | None | $200 | $25 | — | $100 |
| **6** | Crescent Way | Property | Light Blue | $100 | $6 | $50 | $50 |
| **7** | Fortune (Chance) | Normal | None | — | — | — | — |
| **8** | Beacon Street | Property | Light Blue | $100 | $6 | $50 | $50 |
| **9** | Skyline Boulevard | Property | Light Blue | $120 | $8 | $50 | $60 |
| **10** | Prison / Just Visiting | Normal | None | — | — | — | — |
| **11** | Grand Avenue | Property | Pink | $140 | $10 | $100 | $70 |
| **12** | Power Grid | Utility | None | $150 | (Dice $\times 4$) | — | $75 |
| **13** | Highland Park | Property | Pink | $140 | $10 | $100 | $70 |
| **14** | Lakeview Drive | Property | Pink | $160 | $12 | $100 | $80 |
| **15** | North Shore Line | Railway | None | $200 | $25 | — | $100 |
| **16** | Market Street | Property | Orange | $180 | $14 | $100 | $90 |
| **17** | Community Fund | Normal | None | — | — | — | — |
| **18** | Canal Street | Property | Orange | $180 | $14 | $100 | $90 |
| **19** | Union Square | Property | Orange | $200 | $16 | $100 | $100 |
| **20** | Free Parking | Normal | None | — | — | — | — |
| **21** | Pioneer Plaza | Property | Red | $220 | $18 | $150 | $110 |
| **22** | Fortune (Chance) | Normal | None | — | — | — | — |
| **23** | Heritage Way | Property | Red | $220 | $18 | $150 | $110 |
| **24** | Capital Boulevard | Property | Red | $240 | $20 | $150 | $120 |
| **25** | Western Express | Railway | None | $200 | $25 | — | $100 |
| **26** | Sterling Heights | Property | Yellow | $260 | $22 | $150 | $130 |
| **27** | Golden Gate Road | Property | Yellow | $260 | $22 | $150 | $130 |
| **28** | Water Works | Utility | None | $150 | (Dice $\times 4$) | — | $75 |
| **29** | Diamond Parkway | Property | Yellow | $280 | $24 | $150 | $140 |
| **30** | Go To Prison | Normal | None | — | — (Arrest) | — | — |
| **31** | Royal Crescent | Property | Green | $300 | $26 | $200 | $150 |
| **32** | Regal Row | Property | Green | $300 | $26 | $200 | $150 |
| **33** | Community Fund | Normal | None | — | — | — | — |
| **34** | Embassy Gardens | Property | Green | $320 | $28 | $200 | $160 |
| **35** | South Coast Line | Railway | None | $200 | $25 | — | $100 |
| **36** | Fortune (Chance) | Normal | None | — | — | — | — |
| **37** | Summit Ridge | Property | Dark Blue | $350 | $35 | $200 | $175 |
| **38** | Wealth Tax | Normal | None | — | — (Tax: $100) | — | — |
| **39** | Celestial Heights | Property | Dark Blue | $400 | $50 | $200 | $200 |

---

## 2. Monopoly Group Sizes (`state.TotalInColorGroup`)

* **Brown**: 2 tiles (IDs 1, 3)
* **Dark Blue**: 2 tiles (IDs 37, 39)
* **Light Blue**: 3 tiles (IDs 6, 8, 9)
* **Pink**: 3 tiles (IDs 11, 13, 14)
* **Orange**: 3 tiles (IDs 16, 18, 19)
* **Red**: 3 tiles (IDs 21, 23, 24)
* **Yellow**: 3 tiles (IDs 26, 27, 29)
* **Green**: 3 tiles (IDs 31, 32, 34)

---

## 3. Mathematical Formulas

### Railway Rent:
$$\text{Rent} = 25 \times 2^{(\text{RailwaysOwned} - 1)}$$
* 1 station = $25
* 2 stations = $50
* 3 stations = $100
* 4 stations = $200

### Utility Rent:
$$\text{Rent} = \begin{cases} (\text{Dice}_1 + \text{Dice}_2) \times 4 & \text{if 1 utility owned} \\ (\text{Dice}_1 + \text{Dice}_2) \times 10 & \text{if 2 utilities owned} \end{cases}$$

### Demortgage Cost:
$$\text{DemortgageCost} = \text{MortgageValue} + \left\lfloor\frac{\text{MortgageValue}}{10}\right\rfloor$$
*(Principal plus 10% standard bank interest).*

---

## 4. Key Rule Constraints

1. **Even-Building Rule**:
   * A house cannot be built on a property if its current development level exceeds the minimum development level on any other property in the group by more than 0 ($\text{Level} \le \text{MinLevel}$).
   * 4 houses $\to$ 1 hotel (resets `Houses = 0, Hotel = true`).
2. **Even-Selling Rule**:
   * Must sell from the property with the maximum development level in the color group.
   * Bank refunds exactly $50\%$ of the construction cost.
3. **Mortgage Blocking**:
   * Cannot mortgage any property in a color group if *any* property in that group has houses or a hotel.
   * Cannot build on any property in a color group if *any* property in that group is mortgaged.
