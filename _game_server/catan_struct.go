package _game_server

type Resource int

const (
	None Resource = iota
	Sheep
	Rock
	Wheat
	Wood
	Clay
)

type ActionCard int

const (
	Knight ActionCard = iota
	Monopoly
	VictoryPoint
	RoadBuilding
	YearOfPlenty
)

//------------------------------------------------------------------------------------------------------
// Boardgame
/* ID Reference
           ___     ___     ___
          / 0 \___/ 1 \___/ 2 \
          \___/   \___/   \___/
       ___     ___     ___     ___
      / 3 \___/ 4 \___/ 5 \___/ 6 \
      \___/   \___/   \___/   \___/
    ___     ___     ___     ___     ___
   / 7 \___/ 8 \___/ 9 \___/10 \___/11 \
   \___/   \___/   \___/   \___/   \___/
       ___     ___     ___     ___
      /12 \___/13 \___/14 \___/15 \
      \___/   \___/   \___/   \___/
           ___     ___     ___
          /17 \___/18 \___/19 \
          \___/   \___/   \___/
*/

type Hex struct {
	ID           int
	ResourceType Resource
	Value        int // The dice number (2-12)
	HasRobber    bool
}

type Corner struct {
	ID            int
	OwnerID       int // -1 if empty, otherwise Player ID
	IsCity        bool
	AdjacentHexes []int // IDs of hexes that touch this corner
}

type Edge struct {
	ID      int
	OwnerID int
	Corners [2]int // The two corner IDs this edge connects
}

type Board struct {
	Hexes   map[int]*Hex
	Corners map[int]*Corner
	Edges   map[int]*Edge
}

//------------------------------------------------------------------------------------------------------
// structs

type ResourceMap map[Resource]int

type Player struct {
	ID            int
	Resources     ResourceMap
	ActionCards   []ActionCard
	KnightsPlayed int
	LongestRoad   int
	Points        int
}

type Bank struct {
	Resources   ResourceMap
	ActionCards []ActionCard
}

// Global costs for reference
var Costs = map[string]ResourceMap{
	"settlement": {Sheep: 1, Wheat: 1, Wood: 1, Clay: 1},
	"city":       {Wheat: 2, Rock: 3},
	"road":       {Wood: 1, Clay: 1},
	"actionCard": {Sheep: 1, Wheat: 1, Rock: 1},
}

// TODO
var settlement = ResourceMap{Sheep: 1, Wheat: 1, Wood: 1, Clay: 1}
var city = ResourceMap{Wheat: 2, Rock: 3}
var road = ResourceMap{Wood: 1, Clay: 1}
var actionCard = ResourceMap{Sheep: 1, Wheat: 1, Rock: 1}
