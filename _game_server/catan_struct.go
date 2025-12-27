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

type DevCard int

const (
	Knight DevCard = iota
	Monopoly
	VictoryPoint
	RoadBuilding
	YearOfPlenty
)

//------------------------------------------------------------------------------------------------------
// Boardgame

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
	DevCards      []DevCard
	KnightsPlayed int
	LongestRoad   int
	Points        int
}

type Bank struct {
	Resources ResourceMap
	DevDeck   []DevCard
}

// Global costs for reference
var Costs = map[string]ResourceMap{
	"settlement": {Sheep: 1, Wheat: 1, Wood: 1, Clay: 1},
	"city":       {Wheat: 2, Rock: 3},
	"road":       {Wood: 1, Clay: 1},
	"devCard":    {Sheep: 1, Wheat: 1, Rock: 1},
}
