package tile

import (
	"github.com/LandGrab/gamecore/errors"
	"github.com/LandGrab/gamecore/types"
)

const BoardTileCount = 40

var defaultBoard [BoardTileCount]Tile

func init() {
	defaultBoard = [BoardTileCount]Tile{
		// 0: Start (Go)
		NewNormal(0, "Start"),

		// 1: Brown 1 (Mediterranean Ave)
		NewProperty(PropertyConfig{
			ID: 1, Name: "Old Harbor Road", Color: ColorBrown,
			Cost: 60, Rent: 2, Rent1House: 10, Rent2House: 30, Rent3House: 90, Rent4House: 160, RentHotel: 250,
			CostHouse: 50, CostHotel: 50, MortgageValue: 30,
		}),

		// 2: Community Chest 1
		NewNormal(2, "Community Fund"),

		// 3: Brown 2 (Baltic Ave)
		NewProperty(PropertyConfig{
			ID: 3, Name: "Fisherman's Wharf", Color: ColorBrown,
			Cost: 60, Rent: 4, Rent1House: 20, Rent2House: 60, Rent3House: 180, Rent4House: 320, RentHotel: 450,
			CostHouse: 50, CostHotel: 50, MortgageValue: 30,
		}),

		// 4: Income Tax
		NewNormal(4, "Income Tax"),

		// 5: Railway 1 (Reading Railroad)
		NewRailway(5, "Metro Central Station", 200, 25, 100),

		// 6: Light Blue 1 (Oriental Ave)
		NewProperty(PropertyConfig{
			ID: 6, Name: "Crescent Way", Color: ColorLightBlue,
			Cost: 100, Rent: 6, Rent1House: 30, Rent2House: 90, Rent3House: 270, Rent4House: 400, RentHotel: 550,
			CostHouse: 50, CostHotel: 50, MortgageValue: 50,
		}),

		// 7: Chance 1
		NewNormal(7, "Fortune"),

		// 8: Light Blue 2 (Vermont Ave)
		NewProperty(PropertyConfig{
			ID: 8, Name: "Beacon Street", Color: ColorLightBlue,
			Cost: 100, Rent: 6, Rent1House: 30, Rent2House: 90, Rent3House: 270, Rent4House: 400, RentHotel: 550,
			CostHouse: 50, CostHotel: 50, MortgageValue: 50,
		}),

		// 9: Light Blue 3 (Connecticut Ave)
		NewProperty(PropertyConfig{
			ID: 9, Name: "Skyline Boulevard", Color: ColorLightBlue,
			Cost: 120, Rent: 8, Rent1House: 40, Rent2House: 100, Rent3House: 300, Rent4House: 450, RentHotel: 600,
			CostHouse: 50, CostHotel: 50, MortgageValue: 60,
		}),

		// 10: Jail / Just Visiting
		NewNormal(10, "Prison / Just Visiting"),

		// 11: Pink 1 (St. Charles Place)
		NewProperty(PropertyConfig{
			ID: 11, Name: "Grand Avenue", Color: ColorPink,
			Cost: 140, Rent: 10, Rent1House: 50, Rent2House: 150, Rent3House: 450, Rent4House: 625, RentHotel: 750,
			CostHouse: 100, CostHotel: 100, MortgageValue: 70,
		}),

		// 12: Utility 1 (Electric Company)
		NewUtility(12, "Power Grid", 150, 75),

		// 13: Pink 2 (States Ave)
		NewProperty(PropertyConfig{
			ID: 13, Name: "Highland Park", Color: ColorPink,
			Cost: 140, Rent: 10, Rent1House: 50, Rent2House: 150, Rent3House: 450, Rent4House: 625, RentHotel: 750,
			CostHouse: 100, CostHotel: 100, MortgageValue: 70,
		}),

		// 14: Pink 3 (Virginia Ave)
		NewProperty(PropertyConfig{
			ID: 14, Name: "Lakeview Drive", Color: ColorPink,
			Cost: 160, Rent: 12, Rent1House: 60, Rent2House: 180, Rent3House: 500, Rent4House: 700, RentHotel: 900,
			CostHouse: 100, CostHotel: 100, MortgageValue: 80,
		}),

		// 15: Railway 2 (Pennsylvania Railroad)
		NewRailway(15, "North Shore Line", 200, 25, 100),

		// 16: Orange 1 (St. James Place)
		NewProperty(PropertyConfig{
			ID: 16, Name: "Market Street", Color: ColorOrange,
			Cost: 180, Rent: 14, Rent1House: 70, Rent2House: 200, Rent3House: 550, Rent4House: 750, RentHotel: 950,
			CostHouse: 100, CostHotel: 100, MortgageValue: 90,
		}),

		// 17: Community Chest 2
		NewNormal(17, "Community Fund"),

		// 18: Orange 2 (Tennessee Ave)
		NewProperty(PropertyConfig{
			ID: 18, Name: "Canal Street", Color: ColorOrange,
			Cost: 180, Rent: 14, Rent1House: 70, Rent2House: 200, Rent3House: 550, Rent4House: 750, RentHotel: 950,
			CostHouse: 100, CostHotel: 100, MortgageValue: 90,
		}),

		// 19: Orange 3 (New York Ave)
		NewProperty(PropertyConfig{
			ID: 19, Name: "Union Square", Color: ColorOrange,
			Cost: 200, Rent: 16, Rent1House: 80, Rent2House: 220, Rent3House: 600, Rent4House: 800, RentHotel: 1000,
			CostHouse: 100, CostHotel: 100, MortgageValue: 100,
		}),

		// 20: Free Parking
		NewNormal(20, "Free Parking"),

		// 21: Red 1 (Kentucky Ave)
		NewProperty(PropertyConfig{
			ID: 21, Name: "Pioneer Plaza", Color: ColorRed,
			Cost: 220, Rent: 18, Rent1House: 90, Rent2House: 250, Rent3House: 700, Rent4House: 875, RentHotel: 1050,
			CostHouse: 150, CostHotel: 150, MortgageValue: 110,
		}),

		// 22: Chance 2
		NewNormal(22, "Fortune"),

		// 23: Red 2 (Indiana Ave)
		NewProperty(PropertyConfig{
			ID: 23, Name: "Heritage Way", Color: ColorRed,
			Cost: 220, Rent: 18, Rent1House: 90, Rent2House: 250, Rent3House: 700, Rent4House: 875, RentHotel: 1050,
			CostHouse: 150, CostHotel: 150, MortgageValue: 110,
		}),

		// 24: Red 3 (Illinois Ave)
		NewProperty(PropertyConfig{
			ID: 24, Name: "Capital Boulevard", Color: ColorRed,
			Cost: 240, Rent: 20, Rent1House: 100, Rent2House: 300, Rent3House: 750, Rent4House: 925, RentHotel: 1100,
			CostHouse: 150, CostHotel: 150, MortgageValue: 120,
		}),

		// 25: Railway 3 (B. & O. Railroad)
		NewRailway(25, "Western Express", 200, 25, 100),

		// 26: Yellow 1 (Atlantic Ave)
		NewProperty(PropertyConfig{
			ID: 26, Name: "Sterling Heights", Color: ColorYellow,
			Cost: 260, Rent: 22, Rent1House: 110, Rent2House: 330, Rent3House: 800, Rent4House: 975, RentHotel: 1150,
			CostHouse: 150, CostHotel: 150, MortgageValue: 130,
		}),

		// 27: Yellow 2 (Ventnor Ave)
		NewProperty(PropertyConfig{
			ID: 27, Name: "Golden Gate Road", Color: ColorYellow,
			Cost: 260, Rent: 22, Rent1House: 110, Rent2House: 330, Rent3House: 800, Rent4House: 975, RentHotel: 1150,
			CostHouse: 150, CostHotel: 150, MortgageValue: 130,
		}),

		// 28: Utility 2 (Water Works)
		NewUtility(28, "Water Works", 150, 75),

		// 29: Yellow 3 (Marvin Gardens)
		NewProperty(PropertyConfig{
			ID: 29, Name: "Diamond Parkway", Color: ColorYellow,
			Cost: 280, Rent: 24, Rent1House: 120, Rent2House: 360, Rent3House: 850, Rent4House: 1025, RentHotel: 1200,
			CostHouse: 150, CostHotel: 150, MortgageValue: 140,
		}),

		// 30: Go To Jail
		NewNormal(30, "Go To Prison"),

		// 31: Green 1 (Pacific Ave)
		NewProperty(PropertyConfig{
			ID: 31, Name: "Royal Crescent", Color: ColorGreen,
			Cost: 300, Rent: 26, Rent1House: 130, Rent2House: 390, Rent3House: 900, Rent4House: 1100, RentHotel: 1275,
			CostHouse: 200, CostHotel: 200, MortgageValue: 150,
		}),

		// 32: Green 2 (North Carolina Ave)
		NewProperty(PropertyConfig{
			ID: 32, Name: "Regal Row", Color: ColorGreen,
			Cost: 300, Rent: 26, Rent1House: 130, Rent2House: 390, Rent3House: 900, Rent4House: 1100, RentHotel: 1275,
			CostHouse: 200, CostHotel: 200, MortgageValue: 150,
		}),

		// 33: Community Chest 3
		NewNormal(33, "Community Fund"),

		// 34: Green 3 (Pennsylvania Ave)
		NewProperty(PropertyConfig{
			ID: 34, Name: "Embassy Gardens", Color: ColorGreen,
			Cost: 320, Rent: 28, Rent1House: 150, Rent2House: 450, Rent3House: 1000, Rent4House: 1200, RentHotel: 1400,
			CostHouse: 200, CostHotel: 200, MortgageValue: 160,
		}),

		// 35: Railway 4 (Short Line)
		NewRailway(35, "South Coast Line", 200, 25, 100),

		// 36: Chance 3
		NewNormal(36, "Fortune"),

		// 37: Dark Blue 1 (Park Place)
		NewProperty(PropertyConfig{
			ID: 37, Name: "Summit Ridge", Color: ColorDarkBlue,
			Cost: 350, Rent: 35, Rent1House: 175, Rent2House: 500, Rent3House: 1100, Rent4House: 1300, RentHotel: 1500,
			CostHouse: 200, CostHotel: 200, MortgageValue: 175,
		}),

		// 38: Luxury Tax
		NewNormal(38, "Wealth Tax"),

		// 39: Dark Blue 2 (Boardwalk)
		NewProperty(PropertyConfig{
			ID: 39, Name: "Celestial Heights", Color: ColorDarkBlue,
			Cost: 400, Rent: 50, Rent1House: 200, Rent2House: 600, Rent3House: 1400, Rent4House: 1700, RentHotel: 2000,
			CostHouse: 200, CostHotel: 200, MortgageValue: 200,
		}),
	}
}

// TileAt returns the immutable flyweight reference for any position,
// wrapping circularly around the 40-tile board (e.g., position 40 wraps to 0).
func TileAt(position int) Tile {
	idx := (position%BoardTileCount + BoardTileCount) % BoardTileCount
	return defaultBoard[idx]
}

// TileByID looks up and returns the immutable flyweight reference matching the specified ID,
// or ErrTileNotFound if no tile has that ID.
func TileByID(id types.ID) (Tile, error) {
	if id >= BoardTileCount {
		return nil, errors.ErrTileNotFound
	}
	return defaultBoard[id], nil
}

// TileCount returns the total number of tiles on the board (40).
func TileCount() int {
	return BoardTileCount
}

// AllTiles returns a copy slice of all 40 immutable flyweight tile references.
func AllTiles() []Tile {
	tiles := make([]Tile, BoardTileCount)
	copy(tiles, defaultBoard[:])
	return tiles
}
