package _game_server

import "fmt"

func (g *GameSavestate) EndTurn(player Player) {
	CheckLongestRoad(g.Players, g.LongestRoadOwnerID)
	CheckLargestArmy(g.Players, g.CurrentLargestArmy)

	g.UpdateExtraPoints(g.Players)

	fmt.Sprintf("End of turn of Player: %i . Next player is: %i", player.ID, ((player.ID + 1) % (len(g.Players))))
}

// ----------------------------------------------------------------------------------------------------------

func (g *GameSavestate) UpdateExtraPoints(players map[int]*Player) {
	// Minimum requirement is 5 segments
	thresholdRoad := 5
	// Minimum requirement is 3 Knights
	thresholdKnights := 3

	if g.LongestRoadOwnerID != -1 {
		// To take the trophy, you must strictly BEAT the current owner
		thresholdRoad = g.CurrentLongestRoad
	}
	if g.CurrentLargestArmy != -1 {
		// To take the trophy, you must strictly BEAT the current owner
		thresholdKnights = g.CurrentLargestArmy
	}

	for _, p := range players {
		if p.LongestRoad > thresholdRoad {
			// Remove 2 points from old owner
			if oldOwner, ok := players[g.LongestRoadOwnerID]; ok {
				oldOwner.Points -= 2
			}

			// Assign to new owner
			g.LongestRoadOwnerID = p.ID
			g.CurrentLongestRoad = p.LongestRoad
			p.Points += 2
		}
		if p.KnightsPlayed > thresholdKnights {
			// Remove 2 points from old owner
			if oldOwner, ok := players[g.LargestArmyOwnerID]; ok {
				oldOwner.Points -= 2
			}

			// Assign to new owner
			g.LargestArmyOwnerID = p.ID
			g.CurrentLargestArmy = p.KnightsPlayed
			p.Points += 2
		}
	}
}

func (g *GameSavestate) UpdateLongestRoadTrophy() {
	oldOwnerID := g.LongestRoadOwnerID
	bestID := -1
	bestLen := 4 // Must be at least 5 to qualify

	// Find who currently has the longest road on the board
	for _, p := range g.Players {
		if p.LongestRoad > bestLen {
			bestLen = p.LongestRoad
			bestID = p.ID
		}
	}

	// Case A: The trophy moves to a new person (must STRICTLY beat old owner)
	if oldOwnerID != -1 && bestID != oldOwnerID {
		currentOwner := g.Players[oldOwnerID]
		// If the old owner's road was broken and someone else is now longer
		if bestLen > currentOwner.LongestRoad {
			currentOwner.Points -= 2
			g.Players[bestID].Points += 2
			g.LongestRoadOwnerID = bestID
			g.MaxRoadLength = bestLen
		}
	} else if oldOwnerID == -1 && bestID != -1 { // Case B: First time someone hits 5 segments
		g.Players[bestID].Points += 2
		g.LongestRoadOwnerID = bestID
		g.MaxRoadLength = bestLen
	}
}

// -------------------------------------------------------------------------------------------------------

func CheckLargestArmy(players map[int]*Player, currentOwnerID int) int {
	bestID := currentOwnerID
	maxKnights := 2 // Minimum requirement is 3 knights

	if currentOwnerID != -1 {
		maxKnights = players[currentOwnerID].KnightsPlayed
	}

	for id, p := range players {
		if p.KnightsPlayed > maxKnights {
			maxKnights = p.KnightsPlayed
			bestID = id
		}
	}
	return bestID
}

func CheckLongestRoad(players map[int]*Player, currentOwnerID int) int {
	bestID := currentOwnerID
	longRoad := 4 // Minimum requirement is 5 Roads conneceted

	if currentOwnerID != -1 {
		longRoad = players[currentOwnerID].LongestRoad
	}

	for id, p := range players {
		if p.LongestRoad > longRoad {
			longRoad = p.LongestRoad
			bestID = id
		}
	}
	return bestID
}

func (p *Player) GetLongestRoad(b *Board) int {
	maxLen := 0
	playerEdges := b.GetEdgesByPlayer(p.ID)

	for _, edge := range playerEdges {
		// Start a Depth First Search from both ends of every edge
		len1 := dfsRoad(edge.ID, edge.Corners[0], b, p.ID, map[int]bool{edge.ID: true})
		len2 := dfsRoad(edge.ID, edge.Corners[1], b, p.ID, map[int]bool{edge.ID: true})

		if len1 > maxLen {
			maxLen = len1
		}
		if len2 > maxLen {
			maxLen = len2
		}
	}
	return maxLen
}

func (b *Board) GetEdgesByPlayer(playerID int) []*Edge {
	list := []*Edge{}
	for _, e := range b.Edges {
		if e.OwnerID == playerID {
			list = append(list, e)
		}
	}
	return list
}

func dfsRoad(currentEdgeID int, targetCornerID int, b *Board, playerID int, visited map[int]bool) int {
	corner := b.Corners[targetCornerID]

	// Path is broken if an opponent owns this corner
	if corner.OwnerID != -1 && corner.OwnerID != playerID {
		return 1
	}

	bestSubPath := 0
	// Find neighboring edges
	for _, nextEdge := range b.Edges {
		if !visited[nextEdge.ID] && nextEdge.OwnerID == playerID {
			// Check if this edge touches our target corner
			var nextTarget int
			if nextEdge.Corners[0] == targetCornerID {
				nextTarget = nextEdge.Corners[1]
			} else if nextEdge.Corners[1] == targetCornerID {
				nextTarget = nextEdge.Corners[0]
			} else {
				continue
			}

			// Mark visited and recurse
			visited[nextEdge.ID] = true
			pathLen := dfsRoad(nextEdge.ID, nextTarget, b, playerID, visited)
			if pathLen > bestSubPath {
				bestSubPath = pathLen
			}
			// Backtrack
			delete(visited, nextEdge.ID)
		}
	}

	return 1 + bestSubPath
}
