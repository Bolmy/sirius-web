package _game_server

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	GameSaveState *GameSavestate
}

// Helper to send JSON responses
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//--------------------------------------------------------------------------------------------

func (s *Server) handleRollDice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Simulate roll (2d6)
	roll := RollDice()

	// 2. Run your existing logic
	s.GameSaveState.Board.DistributeResources(roll, s.GameSaveState.Players, s.GameSaveState.Bank)

	// 3. Return the result
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"roll":    roll,
		"players": s.GameSaveState.Players,
	})
}

func (s *Server) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Return the entire Game object
	respondWithJSON(w, http.StatusOK, s.GameSaveState)
}

//------------------------------------------------------------------------------------------------

type BuildRequest struct {
	PlayerID int `json:"player_id"`
	CornerID int `json:"corner_id"`
}

func (s *Server) handleBuildSettlement(w http.ResponseWriter, r *http.Request) {
	var req BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Execute the Game logic we built earlier
	err := s.GameSaveState.BuildSettlementOrCity(req.PlayerID, req.CornerID)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Success"})
}
