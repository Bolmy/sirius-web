package _game_server

import "math/rand"
import "time"

// TODO is it necessary ?
/*func initGame(gameID int, playerID []int) GameSavestate {

	gamestate := GameSavestate{
		ID:          gameID,
		Round:       1,
		Board:       NewStandardBoard(),
		Players:     player,
		Bank:        bank,
		ActionCards: generateActionCardDeck(),
	}

	return gamestate
}*/

func NewStandardBoard() *Board {
	board := &Board{
		Hexes:   make(map[int]*Hex),
		Corners: make(map[int]*Corner),
		Edges:   make(map[int]*Edge),
	}

	// Standard Catan Resource Distribution
	resources := []Resource{
		None,
		Wood, Wood, Wood, Wood,
		Wheat, Wheat, Wheat, Wheat,
		Sheep, Sheep, Sheep, Sheep,
		Rock, Rock, Rock,
		Clay, Clay, Clay,
	}

	// Standard Catan Numbers (excluding 7 for Desert)
	numbers := []int{2, 3, 3, 4, 4, 5, 5, 6, 6, 8, 8, 9, 9, 10, 10, 11, 11, 12}

	// Initialize 19 Hexes
	// In a real app, you'd shuffle these slices first
	numIdx := 0
	for i := 0; i < 19; i++ {
		val := 0
		isRobber := false
		if resources[i] == None {
			isRobber = true
		} else {
			val = numbers[numIdx]
			numIdx++
		}

		board.Hexes[i] = &Hex{
			ID:           i,
			ResourceType: resources[i],
			Value:        val,
			HasRobber:    isRobber,
		}
	}
	return board
}

// Mapping logic: Corner ID -> List of Hex IDs
func (b *Board) SetupAdjacency() {
	// Example: Corner 0 is at the top of the board, touching only Hex 0
	b.Corners[0] = &Corner{ID: 0, OwnerID: -1, AdjacentHexes: []int{0}}

	// Example: Corner 10 is in the middle, touching Hexes 0, 1, and 4
	b.Corners[10] = &Corner{ID: 10, OwnerID: -1, AdjacentHexes: []int{0, 1, 4}}

	// In a full implementation, you would loop through 54 corners
	// and assign their 1-3 neighboring Hex IDs.
}

func NewBank() *Bank {

	// Initialize Bank
	bank := &Bank{
		Resources: ResourceMap{
			Sheep: 19,
			Rock:  19,
			Wheat: 19,
			Wood:  19,
			Clay:  19,
		},
		// Initialize and shuffle 25 ActionCards (14 Knights, 5 VPs, etc.)
		ActionCards: generateActionCardDeck(),
	}
	return bank
}

func NewPlayers(playerIDs []int) map[int]*Player {

	// Initialize Players
	players := make(map[int]*Player)
	for _, id := range playerIDs {
		players[id] = &Player{
			ID:            id,
			Resources:     ResourceMap{Sheep: 0, Rock: 0, Wheat: 0, Wood: 0, Clay: 0},
			KnightsPlayed: 0,
			LongestRoad:   0,
			Points:        0,
		}
	}

	return players
}

func generateActionCardDeck() []ActionCard {
	deck := []ActionCard{}
	// Add 14 Knights
	for i := 0; i < 14; i++ {
		deck = append(deck, Knight)
	}
	// 5 Victory Points
	for i := 0; i < 5; i++ {
		deck = append(deck, VictoryPoint)
	}

	// 2 Road Building
	for i := 0; i < 2; i++ {
		deck = append(deck, RoadBuilding)
	}

	// 2 Year of Plenty
	for i := 0; i < 2; i++ {
		deck = append(deck, YearOfPlenty)
	}

	// 2 Monopoly
	for i := 0; i < 2; i++ {
		deck = append(deck, Monopoly)
	}
	// shuffle the deck
	ShuffleActionCardDeck(deck)
	return deck
}

func ShuffleActionCardDeck(deck []ActionCard) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}
