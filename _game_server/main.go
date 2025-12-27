package _game_server

type GameSavestate struct {
	ID int

	Round    int
	Board    *Board
	Bank     *Bank
	Devcards []DevCard
	Players  map[int]*Player

	CurrentLongestRoad int // Tracks the currently longest road length
	LongestRoadOwnerID int // Tracks who currently has the 2 points
	MaxRoadLength      int // The length to beat (min 5)

	CurrentLargestArmy int // Tracks the currently highest amount of knights
	LargestArmyOwnerID int // Tracks who currently has the 2 points
	MaxKnightsPlayed   int // The Knights to beat (min 3)
}

func main() {

}
