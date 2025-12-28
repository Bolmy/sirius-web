package _game_server

import "net/http"

type GameSavestate struct {
	ID int

	Round       int
	Board       *Board
	Bank        *Bank
	ActionCards []ActionCard
	Players     map[int]*Player

	CurrentLongestRoad int // Tracks the currently longest road length
	LongestRoadOwnerID int // Tracks who currently has the 2 points
	MaxRoadLength      int // The length to beat (min 5)

	CurrentLargestArmy int // Tracks the currently highest amount of knights
	LargestArmyOwnerID int // Tracks who currently has the 2 points
	MaxKnightsPlayed   int // The Knights to beat (min 3)
}

func main() {
	// Initialize your game
	game := &GameSavestate{
		Board: NewStandardBoard(),
		Bank:  NewBank(),
	}
	server := &Server{GameSaveState: game}

	// Routes
	http.HandleFunc("/roll", server.handleRollDice)
	http.HandleFunc("/build/settlement", server.handleBuildSettlement)
	http.HandleFunc("/status", server.handleGetStatus)

	println("Catan Backend running on :8080")
	http.ListenAndServe(":8080", nil)
}
